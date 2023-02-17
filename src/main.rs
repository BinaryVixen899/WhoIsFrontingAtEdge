//! Adapted from the Default Compute@Edge template program.

use fastly::http::{header, request, Method, StatusCode};
use fastly::{mime, object_store, Backend, Error, ObjectStore, Request, Response};

#[derive(serde::Deserialize)]
struct Pkdata {
    // id: i32,
    // timestamp: i32,
    members: Vec<Pkmember>,
}

#[derive(serde::Deserialize, serde::Serialize)]
struct Pkmember {
    id: String,
    name: String,
}

#[fastly::main]
fn main(req: Request) -> Result<Response, Error> {
    {}
    // Define backend names
    const PLURALKIT_BACKEND: &str = "pluralkit";

    // Filter request methods...
    match req.get_method() {
        // Allow GET and HEAD requests.
        &Method::GET | &Method::HEAD => (),

        // Deny anything else.
        _ => {
            return Ok(Response::from_status(StatusCode::METHOD_NOT_ALLOWED)
                .with_header(header::ALLOW, "GET, HEAD")
                .with_body_text_plain("This method is not allowed\n"))
        }
    };

    // Pattern match on the path...
    match req.get_path().to_lowercase().as_str() {
        // Make incoming string lowercase

        // If request is to the `/` path...
        "/" => {
            // Send a default synthetic response.
            // Return the "But nothing happefd page"
            Ok(Response::from_status(StatusCode::OK)
                .with_content_type(mime::TEXT_HTML_UTF_8)
                .with_body(include_str!("butnothinghappened.html")))
        }
        "/fronting" => fronting(PLURALKIT_BACKEND, &req),

        "/species" => species(&req),

        // Catch all other requests and return a magikarp 404.
        _ => Ok(Response::from_status(StatusCode::NOT_FOUND)
            .with_content_type(mime::TEXT_HTML_UTF_8)
            // .with_body_text_html(format!("test {alter}", alter = "test").as_str()))
            .with_body(include_str!("magikarp.html"))),
    }
}

fn fronting(backend: &'static str, req: &Request) -> Result<Response, Error> {
    let bereq = Request::get("https://api.pluralkit.me/v2/systems/rzwbg/fronters")
        // sadly pluralkit does not support 304s
        .with_pass(true)
        .with_header("accept", "*/*")
        .with_header("host", "api.pluralkit.me");

    let beresp = bereq.send(backend);

    // Checking for Send Errors
    if let Err(e) = beresp {
        // Log what we encountered
        eprintln!("We've encountered a send error {}:", e);
        // Clearly state it's not our fault
        return Ok(Response::from_status(StatusCode::BAD_GATEWAY));
    }
    // otherwise.. (we should never hit this)
    let mut beresp = beresp.expect("To have handled an invalid send");

    // if we got a successful response
    if !beresp.get_status().is_success() {
        eprintln!(
            "Instead of giving us a 200, PK gave us a {}",
            beresp.get_status()
        );
        return Ok(Response::from_status(StatusCode::SERVICE_UNAVAILABLE));
    }

    // parse the json
    let my_data = beresp.take_body_json::<Pkdata>();
    if let Err(e) = my_data {
        eprintln!("JSON PARSING FAILURE: {}", e);
        return Ok(Response::from_status(StatusCode::SERVICE_UNAVAILABLE));
    }

    // We should never hit this
    let my_data = my_data.expect("To have handled a case in which we had bad JSON data");

    // So now we just need to handle this differently depending on if its someone asking for json or not
    let result = match req.get_header("Accept") {
        Some(x) if x == "application/json" => Some(x),
        _ => None,
    };

    if result.is_some() {
        // Send json of the fronting member back to the client
        Ok(Response::from_status(StatusCode::OK)
            .with_body_json(&my_data.members[0].name)
            .expect("We are able to parse the JSON back")
            .with_content_type(mime::APPLICATION_JSON))
    } else {
        //Display who is fronting in a nice page
        return Ok(Response::from_status(StatusCode::OK)
            .with_content_type(mime::TEXT_HTML_UTF_8)
            .with_body_text_html(
                format!(
                    " <head> <meta charset=\"UTF-8\">
<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">
<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
<title>Document</title>
</head>
<body>
<p>\"The current fronter is: {alter}!\"</p>
</body>
</html>",
                    alter = my_data.members[0].name
                )
                .as_str(),
            ));
    }
}

fn species(req: &Request) -> Result<Response, Error> {
    let foxbox = ObjectStore::open("foxbox").expect("We got an objectstore!");
    let foxbox = foxbox.expect("The objectstore exists!");
    let species = match foxbox.lookup("Species") {
        Ok(s) => match s {
            Some(s) => s.into_string(),

            None => {
                return Ok(Response::from_status(StatusCode::SERVICE_UNAVAILABLE));
            }
        },
        Err(e) => {
            return Ok(Response::from_status(StatusCode::SERVICE_UNAVAILABLE));
        }
    };

    let result = match req.get_header("Accept") {
        Some(x) if x == "application/json" => Some(x),
        _ => None,
    };

    if result.is_some() {
        // Send json of the fronting member back to the client
        Ok(Response::from_status(StatusCode::OK)
            .with_body_json(&species)
            .expect("We are able to parse the JSON back")
            .with_content_type(mime::APPLICATION_JSON))
    } else {
        //Display who is fronting in a nice page
        return Ok(Response::from_status(StatusCode::OK)
            .with_content_type(mime::TEXT_HTML_UTF_8)
            .with_body_text_html(
                format!(
                    " <head> <meta charset=\"UTF-8\">
                    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">
                    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">
                    <title>Document</title>
                    </head>
                    <body>
                    <p>\"Ellen is currently a: {}!\"</p>
                    </body>
                    </html>",
                    species.as_str()
                )
                .as_str(),
            ));
    }
}
