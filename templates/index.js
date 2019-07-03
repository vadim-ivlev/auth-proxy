function loginFormSubmit(event) {
    event.preventDefault()
    $("#result0").html("")
    $("#result1").html("")
    $("#form0").ajaxSubmit({
        url: "/login",
        type: "POST",
        // success: function(response) {$('#result0').text(JSON.stringify(response, null,'  '));}
        success: function(response) {
            $("#result0").jsonViewer(response, {collapsed: true, rootCollapsable: false})
            document.location.reload()
        }
    })
}

function logoutFormSubmit(event) {
    event.preventDefault()
    $("#result0").html("")
    $("#result1").html("")
    $("#form1").ajaxSubmit({
        url: "/logout",
        type: "GET",
        // success: function(response) {$('#result0').text(JSON.stringify(response, null,'  '));}
        success: function(response) {
            $("#result1").jsonViewer(response, {collapsed: true, rootCollapsable: false})
            document.location.reload()
        }
    })
}
