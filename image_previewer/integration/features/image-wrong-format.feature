# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to resize an image
  As an image of wrong format, while has correct extension
  I will receive an appropriate error

  Scenario: then user try to find an image with wrong format, but correct extension, an error has to be received
    When I send GET request to "/fill/50/50/nginx/7z_100x100.jpg" for an image 7z_100x100.jpg
    Then the response code 500 and the response payload "invalid JPEG format" for wrong format image