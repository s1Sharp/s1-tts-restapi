RESOURCE	HTTP    METHOD	ROUTE	        DESCRIPTION

* users	    GET	    /api/v1/users/me	        return the logged-in user’s information
* auth	    POST	/api/v1/auth/register	    Register a new user
* auth	    POST	/api/v1/auth/login	        Login registered user
* auth	    GET	    /api/v1/auth/refresh	    Refresh expired access token
* auth	    GET	    /api/v1/auth/logout	    Logout user

* tasks       GET     /api/v1/tasks/find    
* tasks       GET     /api/v1/tasks/:taskId
* tasks       POST    /api/v1/tasks/

snap install golangci-lint