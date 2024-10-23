# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As an image not found on a remote server
  I will receive an appropriate error

  Scenario: then user try to find a non-existing image on server, an error has to be received
    When I send GET request to "/fill/300/200/nginx/_gopher_fake_1024x504.jpg" for non-existing image _gopher_fake_1024x504.jpg
    Then the response code should be 404
    And the response payload must contain "404 Not Found"