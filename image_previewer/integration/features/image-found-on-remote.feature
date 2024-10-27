# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As an image found on a remote server
  I will receive an appropriate image

  Scenario: then user try to find an existing image on server, an image has to be received
    When I send GET request to "/fill/300/200/nginx/gopher_2000x1000.jpg" for an existing image gopher_2000x1000.jpg
    Then the response code 200 and the response payload must be JPEG image found remotely