// Authentication extension for HTMX
htmx.defineExtension('authentication-header', {
    onEvent: function(name, evt) {
        if (name == "htmx:configureRequest") {
            var authentication = sessionStorage.getItem("Authentication")
            if (authentication != null) {
                evt.detaul.headers["Authentication"] = authentication
            }
        }
    },
    transformResponse: function(text, xhr, elt) {

         // if status is "Unauthorized"
        if (xhr.status == 401) {
            window.location = "/signin"
        }

        var authentication = xhr.getResponseHeader("Authentication")
        if (authentication != null) {
            sessionStorage.setItem("Authentication", authentication)
        }
    }
})