# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As an image found in the local cache
  I will receive an appropriate image

  Scenario: then user try to find an existing image, which is already present in the local cache, an image has to be received
    When I send GET request to "/fill/300/200/nginx/_gopher_original_1024x504.jpg" for an existing image _gopher_original_1024x504.jpg
    Then the response code should be 200
    And the response payload must be a valid JPEG image