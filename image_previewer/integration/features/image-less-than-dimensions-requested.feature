# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As an image is less, than dimensions requested
  I will receive an appropriate error

  Scenario: then user try to find an image, which is less, than dimensions requested, an error has to be received
    When I send GET request to "/fill/500/50/nginx/gopher_50x50.jpg" for an image gopher_50x50.jpg
    Then the response code should be 400
    And the response payload must contain "wrong dimensions"