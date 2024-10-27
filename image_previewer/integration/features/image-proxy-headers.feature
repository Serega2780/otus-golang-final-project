# file: features/event.feature

# http://localhost:8585/
# http://previewer:8585/

Feature: Image previewer
  In order to check if all original headers are proxied to remote host
  As an Authorization header is checked by remote host
  I won't send it in the original request and will receive 401 Unauthorized

  Scenario: then user try to find a remote image by sending a request without Authorization header
    When I send GET request to "/fill/300/200/nginx/gopher_333x666.jpg" for an existing image gopher_333x666.jpg
    Then the response code 401, the response payload "401 Unauthorized", because of auth header absence