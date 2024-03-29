# Blues Wireless

An example of creating a simple Go application on Heroku that is capable of ingesting routed notes.

## Installation

### Setup
1. Visit [Getting Started on Heroku with Go](https://devcenter.heroku.com/articles/getting-started-with-go)
2. Follow the turorial through in its entirety
3. Replace contents of main.go with the main.go code contents found in this repository
4. Update dependencies, commit the new code, and deploy your changes to Heroku
    ```console
    $ go mod tidy
    $ go mod vendor
    $ git add -A .
    $ git commit -m "/db"
    $ git push heroku master
    ```

### Verify installation
  1. From the command line, launch the /route GET handler from your Heroku application (This should return an empty page, as the database table storing forwarded route data is still empty)
        ```console
        $ heroku open route
        ```
  2. Using Postman (or another similar utility), issue a POST request to the URL opened in the last step, with the following JSON in the Body, configured as (application/json)
        ```javascript
        { "Blues Route Test": "Hello from Blues!" }
        ```
  3. Refresh the browser window opened by Step 1 and the following line should be displayed
        ```
        Read from DB: { "Blues Route Test": "Hello from Blues!" }
        ```
### Configure a route in your Notehub project:
  1. Log in to notehub.io and select "Personal Project"
  2. Click on "Routes" in the left navigation
  3. Click "New Route" button in the upper right
  4. Give the Route a name - "My Route"
  5. Ensure the "Route Type" dropdown has a selected value of "General HTTP/HTTPS Request/Response"
  6. Ensure the "Route data from devices" dropdown has a selected value of "All devices"
  7. Ensure the "Route data from notefiles" dropdown has a selected value of "All notefiles"
  8. Ensure the "Transform JSON data before routing" dropdown has a selected value of "No transformation"
  9. Enter the Url from Step 1 of installation verification above into the "Route ro service URL" text box
  10. Ensure the "Route at a maximum rate" dropdown has a selected value of "Unlimited"
  11. Ensure the "Enabled" checkbox in upper right is checked
  12. CLick "Save Route" button in upper right

### Route notefiles to your new route handler
1. Save a notefile to the notecard and sync
2. Refresh the browser window opened by Step 1 in the installation verification
3. The page should update and display "Read from DB:" with the JSON representing the newly added notefile!


