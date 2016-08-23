Glukit
======

Requirements:
=============
  1. Install Google App Engine GO SDK:
  [https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
  2. Read the Getting Started case in case of need:
  [https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)
  3. Get ruby (only for SCSS support), compass, bower and claymate (for [gumby](http://www.gumbyframework.com) and front-end) 
  4. [safekeeper](https://github.com/alexandre-normand/safekeeper) for generating source with client ids/secrets: `go get github.com/alexandre-normand/safekeeper`
  4. Run `./setup.sh` (needs to be done once)
  5. Go in `./view` and run `claymate install`. 
  6. Setup environment variables with the following and run `goapp generate github.com/alexandre-normand/glukit/app/secrets`:
  
    ```
    LOCAL_CLIENT_ID="google-client-id-for-localhost"
    LOCAL_CLIENT_SECRET="google-client-secret-for-localhost"
    PROD_CLIENT_ID=""
    PROD_CLIENT_SECRET=""
    TEST_STRIPE_KEY="stripe-test-key"
    TEST_STRIPE_PUBLISHABLE_KEY="stripe-test-publishable-key"
    PROD_STRIPE_KEY=""
    PROD_STRIPE_PUBLISHABLE_KEY=""
    GLUKLOADER_CLIENT_ID="glukloader-client-id (make one up)"
    GLUKLOADER_CLIENT_SECRET="glukloader-client-secret (make one up)"
    GLUKLOADER_SHARE_EDITION_CLIENT_ID="glukloader-client-id (make one up)"
    GLUKLOADER_SHARE_EDITION_CLIENT_SECRET="glukloader-client-secret (make one up)"
    POSTMAN_CLIENT_ID=""
    POSTMAN_CLIENT_SECRET=""
    SIMPLE_CLIENT_ID=""
    SIMPLE_CLIENT_SECRET=""
    CHROMADEX_CLIENT_ID=""
    CHROMADEX_CLIENT_SECRET=""
    ```

Misc
====
To make `SCSS` changes, use `compass build` or `compass watch`.

Running it:
===========
From <repo path>, execute ```goapp serve``` and hit [http://localhost:8080](http://localhost:8080).
