# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As a remote server does not exist
  I will receive an appropriate error

  Scenario: then user try to find image on non-existing server, an error has to be received
    When I send GET request to "/fill/300/200/raw.githubusercontent2.com/_gopher_original_1024x504.jpg" for non-existing server nginx2
    Then the response code should be 500
    And the response payload must contain "no such host"