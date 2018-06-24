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
  5. Go in `./view` and run `claymate install`. 
  6. Setup environment variables with the following and run `go generate github.com/alexandre-normand/glukit/app/secrets`:
    
    ```
    LOCAL_CLIENT_ID="google-client-id-for-localhost"
    LOCAL_CLIENT_SECRET="google-client-secret-for-localhost"
    PROD_CLIENT_ID=""
    PROD_CLIENT_SECRET=""
    TEST_STRIPE_KEY="stripe-test-key"
    TEST_STRIPE_PUBLISHABLE_KEY="stripe-test-publishable-key"
    PROD_STRIPE_KEY=""
    PROD_STRIPE_PUBLISHABLE_KEY=""
    ```

  7. Setup client ids and secrets using the datastore [administration UI](https://console.cloud.google.com/datastore/entities/query?project=glukit&ns=&kind=osin.client) (as an `osin.client` entity):
  
    * Generate a client secret using something like `openssl rand -base64 24`
    * Generate a client id using something like `echo "`openssl rand -hex 14`.mygluk.it"`
    * Note the `RedirectUri` expected by the authenticating application (i.e. `x-glukloader://oauth/callback`)
    * Create a new `osin.client` entity using those values. The `key.identifier` should match the generated client id.

Misc
====
To make `SCSS` changes, use `compass build` or `compass watch`.

Running it:
===========
From <repo path>, execute ```goapp serve``` and hit [http://localhost:8080](http://localhost:8080).
