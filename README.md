Glukit
======

Requirements:
=============
  1. Install Google App Engine GO SDK:
  [https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
  2. Read the Getting Started case in case of need:
  [https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment)
  3. Get ruby (only for SCSS support), compass, bower and claymate (for [gumby](http://www.gumbyframework.com) and front-end) 
  4. Run `./setup.sh` (needs to be done once)
  5. Go in `./view` and run `claymate install`. 

Misc
====
To make `SCSS` changes, use `compass build` or `compass watch`.

Running it:
===========
From <repo path>, execute ```dev_appserver.py app.yaml``` and hit [http://localhost:8080](http://localhost:8080).
