function onError(e, msg){
    alert(msg, e)
}
// var getResponseHandler = (elementID) => (response) => 
// {
//     $(elementID).jsonViewer(response, {collapsed: true, rootCollapsable: false})
//     document.location.reload()
// }

function showResponseIn(elementID) {
    return function(response){
        $(elementID).jsonViewer(response, {collapsed: true, rootCollapsable: false})
        document.location.reload()
    }
} 



function loginFormSubmit(event) {
    event.preventDefault()
    $("#result0").html("")
    $("#result1").html("")
    $("#form0").ajaxSubmit({
        url: "/login",
        type: "POST",
        success: showResponseIn('#result0'),
        error: onError
    })        
}

function logoutFormSubmit(event) {
    event.preventDefault()
    $("#result0").html("")
    $("#result1").html("")
    $("#form1").ajaxSubmit({
        url: "/logout",
        type: "GET",
        success: showResponseIn('#result1'),
        error: onError
    })
}


