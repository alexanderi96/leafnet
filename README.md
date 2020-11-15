# Leafnet

Project started following this guide: https://astaxie.gitbooks.io/build-web-application-with-golang/en/

Leafnet if a project useful to store information and memories about the family tree

# How To Use?

Via script: bash install.sh

This will generate the binary and set up the database. If you want, you can copy the binary and the public folder into a folder of your choice.

Manually:

    go get github.com/alexanderi96/cicerone

  Change dir to the respective folder and create the db file:
    cat schema.sql | sqlite3 cicerone.db
    run go build
    ./cicerone
  
  Remember to set the environment variable for the session cookie store:
  Under windows: 
     $Env:CICERONE_SESSION_KEY="your super secret key"
  
  open localhost:8081

You can change the port in the config file
