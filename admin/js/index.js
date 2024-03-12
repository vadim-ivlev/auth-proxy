// M O D E L  ******************************************************************************************
var model = {
    checkEmailMessage: `На ваш email выслано сообщение от noreply@rg.ru. 
        Откройте его и перейдите по ссылке 
        для подтверждения вашего электронного адреса.`,

    //---------------------------
    oauth2email: "",
    oauth2name: "",
    oauth2id: "",
    //---------------------------
    // url of GraphQL endpoint = appurl + '/graphql'
    priv_origin: null, 
    set appurl(v){
        this.priv_origin = v
        document.getElementById('appUrl').innerHTML = v //'&#x21E2;&nbsp;'+
        buildSocialIcons(v+"/oauthproviders")
        // document.getElementById('graphqlTestLink').href = 'https://graphql-test.vercel.app/?end_point='+v+'/schema&tab_name=auth-proxy'
    },
    get appurl(){
        return this.priv_origin
    },
    //---------------------------
    // templatesCache keeps loaded templates, not to load them repeatedly
    templatesCache :{},

    //---------------------------
    get logined(){
         return (this._loginedUser != null)
    },
 

    //---------------------------
    urlParams: new URLSearchParams(window.location.search),
    //---------------------------
    _loginedUser: null,
    set loginedUser(v) {
        this._loginedUser = v
        if (v) {
            document.getElementById("userTab").innerText = '✎ '+v.username
            // getAuthRoles(model.loginedUser.username)
        } else {
            document.getElementById("userTab").innerText = ""
        }
        refreshApp()
    },
    get loginedUser() {
        return this._loginedUser
    },

    //---------------------------
    _selfRegAllowed: false,
    set selfRegAllowed(v) {
        this._selfRegAllowed = v
        if (v) {
            showElements("#selfRegButton")
            showElements("#selfRegHelp")
            showElements("#reset-buttons")
        } else {
            hideElements("#selfRegButton")
            hideElements("#selfRegHelp")
            hideElements("#reset-buttons")
        }
    },
    get selfRegAllowed() {
        return this._selfRegAllowed
    },
    
    //---------------------------
    _captchaRequired: false,
    set captchaRequired(v) {
        this._captchaRequired = v
        if (v) {
            showElements("#captcha")
        } else {
            hideElements("#captcha")
        }
    },
    get captchaRequired() {
        return this._captchaRequired
    },

    //---------------------------
    _authRoles: null,
    set authRoles(v) {
        this._authRoles = v
        if (this.isAdmin){
            showElements('#usersTab')
            showElements('#rolesTab')
            showElements('#groupsTab')
            showElements('#graphqlTest')
            showElements('#btnNewApp')
            showElements('#settingsButton')
        } else {
            hideElements('#usersTab')
            hideElements('#rolesTab')
            hideElements('#groupsTab')
            hideElements('#graphqlTest')
            hideElements('#btnNewApp')
            hideElements('#settingsButton')

            // showPage('apps')
        }
        renderTemplateFile('mustache/apps.html', model, '.app-search-results')
        
    },
    get authRoles() {
        return this._authRoles
    },
    get isAdmin() {
        if (!model.authRoles) return false
        return model.authRoles.some(e => (e.rolename == "authadmin" || e.rolename == "auditor"))
    },
    
    //---------------------------
    _user: null,
    set user(v) {
        this._user = v
        renderTemplateFile('mustache/user.html', model, '#userPage')
    },
    get user() {
        return this._user
    },
    
    //---------------------------
    _app: null,
    set app(v) {
        this._app = v
        renderTemplateFile('mustache/app.html', model, '#appPage')
    },
    get app() {
        return this._app
    },
    
    //---------------------------
    _group: null,
    set group(v) {
        this._group = v
        renderTemplateFile('mustache/group.html', model, '#groupPage')
    },
    get group() {
        return this._group
    },
    
    //---------------------------
    _users: null,
    set users(v) {
        this._users = v
        renderTemplateFile('mustache/users.html', model, '.user-search-results')
    },
    get users() {
        return this._users
    },
    
    //---------------------------
    _apps: null,
    set apps(v) {
        this._apps = v
        renderTemplateFile('mustache/apps.html', model, '.app-search-results')
    },
    get apps() {
        return this._apps
    },

    //---------------------------
    _groups: null,
    set groups(v) {
        this._groups = v
        renderTemplateFile('mustache/groups.html', model, '.group-search-results')
    },
    get groups() {
        return this._groups
    },

    //---------------------------
    all_app_options: null,
    all_user_options:null,
    all_group_options:null,
    
    //---------------------------
    _allApps: null,
    set allApps(v) {
        this._allApps = v
        this.all_app_options = createOptions(v, "appname", "description", "url")
    },
    get allApps() {
        return this._allApps
    },

    //---------------------------
    _allUsers: null,
    set allUsers(v) {
        this._allUsers = v
        // this.all_user_options = createOptions(v, "username", "fullname", "email")
        this.all_user_options = createOptions(v, "username", "fullname")
    },
    get allUsers() {
        return this._allUsers
    },

    //---------------------------
    _allGroups: null,
    set allGroups(v) {
        this._allGroups = v
        console.debug('allGroups',v)
        this.all_groups_options = createOptions(v, "groupname", "description")
    },
    get allGroups() {
        return this._allGroups
    },

    
    //---------------------------
    _app_user_roles: null,
    set app_user_roles(v) {
        this._app_user_roles = v
    },
    get app_user_roles() {
        return this._app_user_roles
    },

    //---------------------------
    _app_group_roles: null,
    set app_group_roles(v) {
        this._app_group_roles = v
        renderTemplateFile('mustache/app-group-roles.html', model, '.app-group-roles-results')
    },
    get app_group_roles() {
        return this._app_group_roles
    },

    //---------------------------
    _user_group_roles: null,
    set user_group_roles(v) {
        this._user_group_roles = v
        renderTemplateFile('mustache/user-group-roles.html', model, '.user-group-roles-results')
    },
    get user_group_roles() {
        return this._user_group_roles
    },

    //---------------------------
    _params: null,
    set params(v) {
        this._params = v
        renderTemplateFile('mustache/params.html', model, '#paramsPage')
        document.getElementById('appName').innerText = v ? v.app_name : '';
    },
    get params() {
        return this._params
    },

    //---------------------------
    _appstat: null,
    set appstat(v) {
        this._appstat = v
        if (! document.getElementById('gauges')) return

        document.getElementById('divSys').innerText = v.sys +' Mb'
        document.getElementById('divAlloc').innerText = 'allocated: '+v.alloc +' Mb'
        document.getElementById('divTotalAlloc').innerText = 'total allocated: '+v.total_alloc +' Mb'


        drawGauge("req / day", v.requests_per_day, 0,  "divDay")        // 10000
        drawGauge("req / hour", v.requests_per_hour, 0, "divHour")       // 1000
        drawGauge("req / min", v.requests_per_minute, 0, "divMinute")     // 100
        drawGauge("req / sec", v.requests_per_second, 0, "divSecond")      // 10

        redrawMemoryChart()
    },
    get appstat() {
        return this._appstat
    },

 
    
}



// F U N C T I O N S  *********************************************************************************

// function getListOptionProperty(inputId,propertyName) {
//     var input = document.getElementById(inputId)
//     if (!input) return
//     var value = input.value
//     var selector = `#${inputId}>option[value='${value}']`
//     if (!selector) return
//     var option = document.querySelector(selector)
//     if (!option) return
//     var value = option[propertyName]
//     return value
// }

// prevent default
function pd(event) {
    if (event) event.preventDefault()
}


function createOptions(selectValues, valueProp, textProp1, textProp2) {
    var output = []
    selectValues && selectValues.forEach(function(value)
    {
      output.push(`<option id="${value.id}" value="${value[valueProp]}">${value[textProp1]} &nbsp;&nbsp;&nbsp; ${value[textProp2]?value[textProp2]:''}</option>`);
    })
    let optionText = output.join('')
    return optionText
}


function highlightTab(tabid) {
        removeClass('.tab', "underlined")
        var tabid0 = tabid.split("/")[0]
        addClass('#'+tabid0+'Tab', "underlined")   
}


function showPage(pageid, dontpush){
    stopGettingAppstat()

    //распарсить pageidExtended
    var a = pageid.split("/")
    var pageName = a[0]
    var id = a[1]

    highlightTab(pageid)
    
    hideElements('.page')
    showElements('#'+pageName+'Page')


    // setting focus
    var text = document.querySelector('#'+pageName+'Page input[type="text"]')
    if(text) 
        text.focus()


    if (!dontpush){
        if (!history.state || history.state.pageid != pageid ){
            history.pushState({pageid:pageid},pageid, "#"+pageid) 
        }
    }

    if (id) {
        if (pageName == "app"){
            getApp(id)
        } 
        if (pageName == "user"){
            getUser(id)
        }
        if (pageName == "group"){
            getGroup(id)
        }
    }

    if (pageName == 'params')
        startGettingAppstat()

    return false
}


function renderTemplateFile(templateFile, data, targetSelector) {
    // извлекаем шаблон из кэша
    var cachedTemlpate = model.templatesCache[templateFile]
    
    if (cachedTemlpate) {
        renderTemplate(cachedTemlpate) 
        console.debug("from cache:",templateFile)
        return
    }
    
    fetch(templateFile).then(x => x.text()).then( template => {
        // запоминаем шаблон в кэше
        model.templatesCache[templateFile]=template 
        renderTemplate(template)
        console.debug("loaded template:",templateFile)
    })

    return

    function renderTemplate(template) {
        var rendered = Mustache.render(template, data)
        document.querySelector(targetSelector).innerHTML = rendered
    }    

}


var delayTimeout
function delayFunc(f, delay=500) {
   clearTimeout(delayTimeout) 
   delayTimeout = setTimeout(f, delay)  
   return false 
}


function searchApps() {
    if (document.querySelector('#chkLocalSearch').checked) {
        return delayFunc(searchAppsInModel, 100)
    } else {
        return delayFunc(formListAppSubmit)
    }
}

function searchAppsInModel() {
    if (!model.allApps) return
    var text = document.querySelector("#formListApp input[name='search']").value.trim().replace(' ','.*')
    var r = new RegExp(text, 'i')
    let found = model.allApps.filter((v)=>{
        var s = Object.values(v).join(' ')
        return r.test(s)
    })
    model.apps = found
    return false   
}

/**
 * Сортирует массив приложений
 * 
 * @param {string} prop имя поля
 * @param {boolean} asStings сортировать лексикографически
 * @param {boolean} reverse сортировать в обратном порядке
 */
function sortAppsBy(prop, asStings, reverse) {
    if (asStings){
        model._allApps?.sort( (a,b) => (a[prop]+a.appname)>(b[prop]+b.appname)? 1: -1)
        model._apps?.sort( (a,b) => (a[prop]+a.appname)>(b[prop]+b.appname)? 1: -1 )
    }else {
        model._allApps?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1)
        model._apps?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1 )
    }
    if (reverse){
        model._allApps?.reverse()
        model._apps?.reverse()
     }
    model.apps = model._apps
    return false
}


function searchUsers() {
    if (document.querySelector('#chkLocalSearch').checked) {
        return delayFunc(searchUsersInModel, 100)
    } else {
        return delayFunc(formListUserSubmit)
    }
}


function searchUsersInModel() {
    if (!model.allUsers) return
    var text = document.querySelector("#formListUser input[name='search']").value.trim().replace(' ','.*')
    var r = new RegExp(text, 'i')
    let found = model.allUsers.filter( (v) => r.test (Object.values(v).join(' ')) )
    model.users = found
    return false   
}

/**
 * Сортирует массив пользователей 
 * 
 * @param {string} prop имя поля
 * @param {boolean} asStings сортировать лексикографически
 * @param {boolean} reverse сортировать в обратном порядке
 * @returns 
 */
function sortUsersBy(prop, asStings, reverse) {
    if (asStings){
            model._allUsers?.sort( (a,b) => (a[prop]+a.username)>(b[prop]+b.username)? 1: -1)
            model._users?.sort( (a,b) => (a[prop]+a.username)>(b[prop]+b.username)? 1: -1 )
    } else {
            model._allUsers?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1)
            model._users?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1 )
    }
    if (reverse){
        model._allUsers?.reverse()
        model._users?.reverse()
     }
    model.users = model._users
    return false
}

function searchGroups() {
    if (document.querySelector('#chkLocalSearch').checked) {
        return delayFunc(searchGroupsInModel, 100)
    } else {
        return delayFunc(formListGroupSubmit)
    }
}



function searchGroupsInModel() {
    if (!model.allGroups) return
    var text = document.querySelector("#formListGroup input[name='search']").value.trim().replace(' ','.*')
    var r = new RegExp(text, 'i')
    let found = model.allGroups.filter( (v) => r.test (Object.values(v).join(' ')) )
    model.groups = found
    return false   
}

/**
 * Сортирует массив групп
 * 
 * @param {string} prop имя поля
 * @param {boolean} asStings сортировать лексикографически
 * @param {boolean} reverse сортировать в обратном порядке
 */
function sortGroupsBy(prop, asStings, reverse) {
    if (asStings){
        model._allGroups?.sort( (a,b) => (a[prop]+a.groupname)>(b[prop]+b.groupname)? 1: -1)
        model._groups?.sort( (a,b) => (a[prop]+a.groupname)>(b[prop]+b.groupname)? 1: -1 )
    } else {
        model._allGroups?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1)
        model._groups?.sort( (a,b) => (a[prop])>(b[prop])? 1: -1 )
    }
    if (reverse){
        model._allGroups?.reverse()
        model._groups?.reverse()
     }
    model.groups = model._groups
    return false
}


function errorMessage(errorElementID, errMsg) {
    console.debug(errMsg)
    if (errorElementID) {
        document.getElementById(errorElementID).innerText = errMsg
    }
    if (errMsg=="email not confirmed") {
        alert(model.checkEmailMessage)
    }
}


// R E Q U E S T S  *******************************************************

function doGraphQLRequest(query, responseHandler, errorElementID="loginError", xreqidHeader=null, onError) {
    fetch(model.appurl+'/graphql', { 
        // headers: new Headers({cache: "no-cache"}),
        headers: new Headers({"x-req-id": xreqidHeader}),
        method: 'POST', 
        credentials: 'include', 
        body: JSON.stringify({ query: query, variables: {} }) 
    })
        .then(res => { 
            if (res.ok) {
                return res.json();
            }
            new Error(res)
        })
        .then((res) => {
            console.debug("doGraphQLRequest() result=",res)
            if (res.errors){
                let errMsg = res.errors[0].message
                errorMessage(errorElementID, errMsg)
                if (onError && typeof onError === 'function'){ 
                    onError(res)
                }
                return
            }
            if ( responseHandler && typeof responseHandler === 'function') { 
                responseHandler(res)
            }
        })
        .catch( e =>  errorMessage(errorElementID, "doGraphQLRequest: " +e ))    
    
}


// L O G I N  *****************************************************************************

function loginGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
    
    let username = document.getElementById("loginUsername").value
    let password = document.getElementById("loginPassword").value
    let captcha =  document.getElementById("loginCaptcha").value
    let pin =  document.getElementById("loginPin").value
    if (username == "") {
        errorMessage("loginError", 'Заполните поле email')
        return
    }

    isPinRequired(username)
    

    let query =`
    query {
        login(
        username: "${username}",
        password: "${password}",
        captcha: "${captcha}",
        pin: "${pin}"
        )
        }    
    `
    function onSuccess(res){
        getLoginedUser()
    }   
    
    doGraphQLRequest(query, onSuccess, "loginError")
    clearLoginForm()
    return false       
}



function logoutGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
    var query =`
    query {
        logout {
            message
            email
          }
        }
    `

    function onSuccess(res){
        model.loginedUser = null
        model.oauth2email =""
        model.oauth2id =""
        model.oauth2name =""

        refreshApp()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}



function isSelfRegAllowed(event) {
    if (event) event.preventDefault()
    var query =` query { is_selfreg_allowed }`

    function onSuccess(res){
        model.selfRegAllowed = res.data.is_selfreg_allowed
    }
       
    doGraphQLRequest(query, onSuccess)
    return false       
}


function checkUserRequirements(event) {
    if (event) event.preventDefault()
    let username = document.getElementById("loginUsername").value   
    isCaptchaRequired(username)
    isPinRequired(username)
    return false 
}


function isCaptchaRequired(username) {
    var query =`  query { is_captcha_required(  username: "${username}" ) 
        {
            is_required 
            path
        } 
    } `

    function onSuccess(res){
        model.captchaRequired = res.data.is_captcha_required.is_required
        if (model.captchaRequired) {
            getNewCaptcha()
            console.debug("Captcha IS required")
        }
    }
       
    doGraphQLRequest(query, onSuccess)   
}


function isPinRequired(username) {
    errorMessage("loginError", "")
    
    var query =`  query { is_pin_required(  username: "${username}" ) 
        {
            use_pin 
            pinrequired
        } 
    } `

    function onSuccess(res){
        let r = res.data.is_pin_required;
        console.debug("isPinRequired()->", r);
        (r.use_pin && r.pinrequired) ? showElements(".pinclass") : hideElements(".pinclass"); 
    }
    function onError(res){
        console.debug("isPinRequired().onError(res)->", res);
        hideElements(".pinclass")
    }
    doGraphQLRequest(query, onSuccess, undefined, null, onError)   
}



function resetPasswordRest(event) {
    if (event) event.preventDefault()
    errorMessage("resetError", "")
    let username = document.getElementById("loginUsername").value
    if (!username) {
        errorMessage("resetError", 'Заполните поле email')
        return false
    }
    let adminurl = location.origin + location.pathname
    let url = `${model.priv_origin}/reset_password?username=${username}&adminurl=${adminurl}&authurl=${model.priv_origin}`
    fetch(url).then( r => r.json()).then(onSuccess).catch(onError)
    
    function onSuccess(res){
        errorMessage("resetError", res.error)
        if (res.result) {
             alert(res.result)
        } 
    }

    function onError(err){
        errorMessage("resetError", "resetPasswordRest " + err)
    }

    return false 
}



function resetAuthenticator(event) {
    if (event) event.preventDefault()
    errorMessage("resetError", "")
    let username = document.getElementById("loginUsername").value
    if (!username) {
        errorMessage("resetError", 'Заполните поле email')
        return false
    }
    let adminurl = location.origin + location.pathname
    let url = `${model.priv_origin}/reset_authenticator?username=${username}&adminurl=${adminurl}&authurl=${model.priv_origin}`
    fetch(url).then( r => r.json()).then(onSuccess).catch(onError)
    
    function onSuccess(res){
        errorMessage("resetError", res.error)
        if (res.result) {
             alert(res.result)
        } 
    }

    function onError(err){
        errorMessage("resetError", "resetAuthenticator "+err)
    }

    return false 
}



// L O G I N E D   U S E R   **********************************************************************************************************************


function getLoginedUser() {
    model.loginedUser = null
    var query =`
    query {
        get_logined_user {
            description
            email
            fullname
            username
            disabled
            id
            pinhash
            pinrequired
            pinhash_temp  
            emailhash
            emailconfirmed    
        }
    }
    `
    
    function onSuccess(res){
        model.loginedUser = res.data.get_logined_user
        getAuthRoles(model.loginedUser.username)
        getUser(model.loginedUser.username)
    } 
    
    doGraphQLRequest(query, onSuccess)
    return false       
}


function getAuthRoles(username) {
    model.authRoles = null

    var query =`
    query {
        list_app_user_role(
            appname: "auth",
            username: "${username}"
            ) {
                rolename
            }
        }
        `

    function onSuccess(res){
        model.authRoles = res.data.list_app_user_role
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}

// P A R A M S  ***********************************************************************************************************************


function getParams() {
    model.params = null
    var query =`
    query {
        get_params {
            app_name
            max_attempts
            reset_time
            selfreg
            use_captcha
            use_pin
            login_not_confirmed_email
            no_schema
          }
        }
    `

    function onSuccess(res){
        var p = res.data.get_params
        model.params = p
        getAppstatRest()
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}

function setParams(event) {
    if (event) event.preventDefault()
    let selfreg =                    document.querySelector("#formParams input[name='selfreg']").checked
    let use_captcha =                document.querySelector("#formParams input[name='use_captcha']").checked
    let use_pin =                    document.querySelector("#formParams input[name='use_pin']").checked
    let login_not_confirmed_email =  document.querySelector("#formParams input[name='login_not_confirmed_email']").checked
    let no_schema =                  document.querySelector("#formParams input[name='no_schema']").checked
    let max_attempts =               document.querySelector("#formParams input[name='max_attempts']").value
    let reset_time =                 document.querySelector("#formParams input[name='reset_time']").value
    
    var query =`
    query {
        set_params(
            selfreg:                   ${selfreg},
            use_captcha:               ${use_captcha},
            use_pin:                   ${use_pin},
            login_not_confirmed_email: ${login_not_confirmed_email},
            no_schema:                 ${no_schema},
            max_attempts:              ${max_attempts},
            reset_time:                ${reset_time}
        ) {
            max_attempts
            reset_time
            selfreg
            use_captcha
            use_pin
            login_not_confirmed_email
            no_schema
          }
        }
    `

    function onSuccess(res){
        model.params = res.data.set_params
        alert("params saving success")
    }

    doGraphQLRequest(query, onSuccess, "paramsError")
    return false       
}


// function getAppstat(event) {
//     if (event) event.preventDefault()
//     var query =`
//     query {
//         get_stat {
//             alloc
//             requests_per_day
//             requests_per_hour
//             requests_per_minute
//             requests_per_second
//             sys
//             total_alloc
//           }
//         }
//     `

//     function onSuccess(res){
//         var m = res.data.get_stat
//         model.appstat = m
//     } 

//     doGraphQLRequest(query, onSuccess, "appstatError")
//     return false       
// }

function getAppstatRest(event) {
    if (event) event.preventDefault()
    getAppstatRest.counter = getAppstatRest.counter === undefined ? 0: getAppstatRest.counter+1
    console.log("getAppstatRest.counter = ",getAppstatRest.counter)


    fetch(model.appurl + "/stat").then(x => x.json())
    .then( onSuccess )
    .catch( err => {
        console.log("fetch error:",err)
        return
        }
    )  

    function onSuccess(res){
        var m = res
        model.appstat = m
    } 
    return false       
}

model.statInterval = null
function startGettingAppstat(){
    clearInterval(model.statInterval)
    getAppstatRest()
    model.statInterval = setInterval(getAppstatRest, getRefreshInterval())
    console.log('startGettingAppstat')
}

function stopGettingAppstat(){
    clearInterval(model.statInterval)
    console.log('stopGettingAppstat')
}

function setRefreshInterval(event) {
    if (event) event.preventDefault()
    el = document.querySelector("#sliderValue")
    el.innerText = getRefreshInterval()   
    stopGettingAppstat()
    startGettingAppstat()
    return false       
}

function getRefreshInterval() {
    return document.querySelector("#slider")?.value || 3000
}

function initMemoryChart() {

}

function initMemoryDataTable() {
    
}

function redrawMemoryChart() {
     
}


// U S E R S  ***********************************************************************************************************************


function formListUserSubmit(event) {
    if (event) event.preventDefault()
    model.users = null
    let search = document.querySelector("#formListUser input[name='search']").value
    
    var query =`
    query {
        list_user(
        search: "${search}",
        order: "fullname ASC"
        ) {
            length
            list {
              description
              email
              fullname
              username
              disabled
              id
              pinhash
              pinrequired
              pinhash_temp   
              emailhash
              emailconfirmed    
              }
          }
        }        
    `
    
    function onSuccess(res){
        model.users = res.data.list_user.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}




function updateUser(event) {
    if (event) event.preventDefault()
    let id =  document.querySelector("#user-id").innerText
    // let old_username =  document.querySelector("#formUser input[name='old_username']").value
    // let username =      document.querySelector("#formUser input[name='username']").value
    let password =      document.querySelector("#formUser input[name='password']").value
    // let email    =      document.querySelector("#formUser input[name='email']").value
    let fullname =      document.querySelector("#formUser input[name='fullname']").value
    let description =   document.querySelector("#formUser *[name='description']").value
    let disabled =      document.querySelector("#formUser input[name='disabled']").value
    let pinrequired =   document.querySelector("#formUser input[name='pinrequired']").checked
    let emailconfirmed =   document.querySelector("#formUser input[name='emailconfirmed']").checked
    
    var query =`
    mutation {
        update_user(
        id: ${id}
        password: "${password}",
        fullname: "${fullname}",
        description: "${description}",
        disabled: ${disabled},
        pinrequired: ${pinrequired},
        emailconfirmed: ${emailconfirmed},
        ) {
            description
            email
            fullname
            username
            disabled
            id
            pinhash
            pinrequired
            pinhash_temp      
            emailhash
            emailconfirmed    
          }

        }
    `

    function onSuccess(res){
        alert("update_user success")
        refreshData()
        model.user = res.data["update_user"]
        getUser(model.user.username)
    }

    doGraphQLRequest(query, onSuccess, "userError")
    return false       
}

// https://stackoverflow.com/questions/18338890/are-there-any-sha-256-javascript-implementations-that-are-generally-considered-t/48161723#48161723
async function sha256(message) {
    // encode as UTF-8
    const msgBuffer = new TextEncoder().encode(message);                    

    // hash the message
    const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);

    // convert ArrayBuffer to Array
    const hashArray = Array.from(new Uint8Array(hashBuffer));

    // convert bytes to hex string                  
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    return hashHex;
}

async function createUser(event) {
    if (event) event.preventDefault()
    // let username =      document.querySelector("#formUser input[name='username']").value
    let email    =         document.querySelector("#formUser input[name='email']").value
    let password =         document.querySelector("#formUser input[name='password']").value
    let fullname =         document.querySelector("#formUser input[name='fullname']").value
    let description =      document.querySelector("#formUser *[name='description']").value
    let disabled =         document.querySelector("#formUser input[name='disabled']").value
    let pinrequired =      document.querySelector("#formUser input[name='pinrequired']").checked
    let emailconfirmed =   document.querySelector("#formUser input[name='emailconfirmed']").checked
    let noemail =          document.querySelector("#formUser input[name='noemail']").checked
    let addgroup =         document.querySelector("#formUser select[name='addgroup']").value
    let addgroupParam = addgroup ? `addgroup: ${addgroup}` : ""
    
    var query =`
    mutation {
        create_user(
        email: "${email}",
        password: "${password}",
        fullname: "${fullname}",
        description: "${description}",
        disabled: ${disabled},        
        pinrequired: ${pinrequired},
        emailconfirmed: ${emailconfirmed},
        noemail: ${noemail},
        ${addgroupParam}
        ) {
            description
            email
            fullname
            username
            disabled
            id
            pinhash
            pinrequired 
            pinhash_temp       
            emailhash
            emailconfirmed    
          }

        }
    `

    function onSuccess(res){
        alert("create_user success")
        refreshData()
        model.user = res.data["create_user"]
        if (!model.logined) {
            alert(`"${model.user.username}" is created.` + "\n\n" + model.checkEmailMessage )
            showPage('login',true)
        } else {
            getUser(model.user.username)
        }
    }

    // вычисляем секретный заголовок
    const xreqid = await sha256(fullname + email + password)
    doGraphQLRequest(query, onSuccess, "userError", xreqid)
    return false       
}


function getUser(username) {
    model.user = null
    var query =`
    query {
        get_user(
        username: "${username}"
        ) {
            description
            email
            fullname
            username
            disabled
            id
            pinhash
            pinrequired
            pinhash_temp     
            emailhash
            emailconfirmed    
          }
        
        list_app_user_role(
        username: "${username}"
        ) {
            app_description
            appname
            rolename
            username
          }
    }
    `

    function onSuccess(res){
        var u = res.data.get_user
        u.apps = appsOfTheUser(res.data.list_app_user_role)
        // render 
        model.user = u
        fetchGroupsOfTheUser(u.id)
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}



function appsOfTheUser(list_app_user_role) {
    if (!list_app_user_role) return []
    
    let gr = {}
    for (let aur of list_app_user_role ){
        if (!gr[aur.appname]) gr[aur.appname] =[]
        gr[aur.appname].push(aur)
    }

    let arr = []

    for (let [key, value] of Object.entries(gr)) {
        let rec = {}
        rec.appname =key
        rec.app_description = value[0].app_description
        rec.items = value
        arr.push(rec)
    }
    return arr
}



function deleteUser(id) {
    if (!confirm(`Удалить пользователя ${id}?`)) return false

    var query =`
    mutation {
        delete_user(
        id: ${id}
        ) {
            username
          }
        }
    `
    function onSuccess(res){
        model.user = null
        refreshData()
        showPage('users')
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}

function sendPasswordEmail(email) {
    var query =`
    mutation {
        send_password_email( email_to: "${email}") 
        }
    `
    function onSuccess(res){
        alert(JSON.stringify(res, null, 2))
    }   
    function onError(res){
        alert(JSON.stringify(res, null, 2))
    }
    doGraphQLRequest(query, onSuccess, undefined, null, onError)
    return false       
}

function sendAuthenticatorEmail(email) {
    var query =`
    mutation {
        send_authenticator_email( email_to: "${email}") 
        }
    `
    function onSuccess(res){
        alert(JSON.stringify(res, null, 2))
    } 
    function onError(res){
        alert(JSON.stringify(res, null, 2))
    }
    doGraphQLRequest(query, onSuccess, undefined, null, onError)
    return false       
}


function fetchGroupsOfTheUser(user_id) {
    var query =`
    query {
        list_group_user_role(
        user_id: ${user_id}
        ) {
            group_description
            group_groupname
            group_id
            rolename
            user_description
            user_disabled
            user_email
            user_fullname
            user_id
            }        
        }
    `

    function onSuccess(res){
        var u = model.user
        u.groups = groupsOfTheUser(res.data.list_group_user_role)
        // render 
        model.user = u
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function groupsOfTheUser(list_group_user_role) {
    let groups = groupByField(list_group_user_role, 'group_id')

    // преобразуем хэш в массив для отображения в mustache
    let arr = []

    for (let [key, value] of Object.entries(groups)) {
        let rec = {}
        rec.group_id =key
        rec.user_id = value[0].user_id
        rec.user_email = value[0].user_email
        rec.group_groupname = value[0].group_groupname
        rec.group_description = value[0].group_description
        rec.items = value
        arr.push(rec)
    }
    return arr
}



// A P P S  *******************************************************************

function formListAppSubmit(event) {
    if (event) event.preventDefault()
    model.apps = null
    let search = document.querySelector("#formListApp input[name='search']").value
    var query =`
    query {
        list_app(
        search: "${search}",
        order: "description ASC"
        ) {
            length
            list {
              id  
              appname
              description
              url
              rebase
              public
              sign
              xtoken
            }
          }
        }    `

    function onSuccess(res){
        model.apps = res.data.list_app.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}



function createApp(event) {
    if (event) event.preventDefault()
    let appname =     document.querySelector("#formApp input[name='appname']"    ).value
    let url =         document.querySelector("#formApp input[name='url']"        ).value
    let description = document.querySelector("#formApp input[name='description']").value
    let rebase =      document.querySelector("#formApp input[name='rebase']"     ).value
    let _public =     document.querySelector("#formApp input[name='public']"     ).value
    let sign =        document.querySelector("#formApp input[name='sign']"       ).value
    let xtoken =        document.querySelector("#formApp input[name='xtoken']"       ).value
    
    var query =`
    mutation {
        create_app(
        appname: "${appname}",
        url: "${url}",
        description: "${description}",
        rebase: "${rebase}",
        public: "${_public}",
        sign: "${sign}"
        xtoken: "${xtoken}"
        ) {
            description
            appname
            url
            rebase
            public
            sign
            xtoken
          }

        }
    `
    function onSuccess(res){
        alert("create_app success")
        refreshData()
        model.app = res.data["create_app"]
        getApp(appname)
    } 

    doGraphQLRequest(query, onSuccess, "appError")
    return false       
}


function updateApp(event, appOperationName = 'create_app') {
    if (event) event.preventDefault()
    let app_id =      document.querySelector("#formApp input[name='app_id']"    ).value
    // let old_appname = document.querySelector("#formApp input[name='old_appname']"    ).value
    let appname =     document.querySelector("#formApp input[name='appname']"    ).value
    let url =         document.querySelector("#formApp input[name='url']"        ).value
    let description = document.querySelector("#formApp input[name='description']").value
    let rebase =      document.querySelector("#formApp input[name='rebase']"     ).value
    let _public =     document.querySelector("#formApp input[name='public']"     ).value
    let sign =        document.querySelector("#formApp input[name='sign']"       ).value
    let xtoken =        document.querySelector("#formApp input[name='xtoken']"       ).value
    
    // old_appname: "${old_appname}",
    var query =`
    mutation {
        update_app(
        id: ${app_id},
        appname: "${appname}",
        url: "${url}",
        description: "${description}",
        rebase: "${rebase}",
        public: "${_public}",
        sign: "${sign}"
        xtoken: "${xtoken}"
        ) {
            description
            appname
            url
            rebase
            public
            sign
            xtoken
          }

        }
    `
    function onSuccess(res){
        alert("update_app success")
        refreshData()
        model.app = res.data["update_app"]
        getApp(appname)
    } 

    doGraphQLRequest(query, onSuccess, "appError")
    return false       
}



function getApp(appname) {
    model.app = null

    var query =`
    query {
        get_app(
        appname: "${appname}"
        ) {
            appname
            description
            id
            public
            rebase
            sign
            url
            xtoken
          }
        
        list_app_user_role(
        appname: "${appname}"
        ) {
            appname
            rolename
            user_fullname
            user_disabled
            username
          }
        
        }

    `

    function onSuccess(res){
        var a = res.data.get_app
        a.users = usersOfTheApp(res.data.list_app_user_role)
        // render
        model.app = a
        fetchGroupsOfTheApp(a.id)
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function usersOfTheApp(list_app_user_role) {
    let gr = {}
    for (let aur of list_app_user_role ){
        if (!gr[aur.username]) gr[aur.username] =[]
        gr[aur.username].push(aur)
    }

    let arr = []

    for (let [key, value] of Object.entries(gr)) {
        let rec = {}
        rec.username =key
        rec.user_fullname = value[0].user_fullname
        rec.items = value
        arr.push(rec)
    }
    return arr
}


function deleteApp(id) {
    if (!confirm(`Удалить приложение ${id}?`)) return false

    var query =`
    mutation {
        delete_app(
        id: ${id}
        ) {
            appname
          }
        }
    `
    function onSuccess(res){
        model.app = null
        refreshData()
        showPage('apps') ;
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function fetchGroupsOfTheApp(app_id) {
    var query =`
    query {
        list_group_app_role(
        app_id: ${app_id}
        ) {
            app_appname
            app_description
            app_id
            app_url
            group_description
            group_groupname
            group_id
            rolename
            }        
        }
    `

    function onSuccess(res){
        var a = model.app
        a.groups = groupsOfTheApp(res.data.list_group_app_role)
        // render 
        model.app = a
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function groupsOfTheApp(list_group_app_role) {
    let groups = groupByField(list_group_app_role, 'group_id')

    // преобразуем хэш в массив для отображения в mustache
    let arr = []

    for (let [key, value] of Object.entries(groups)) {
        let rec = {}
        rec.group_id =key
        rec.app_id = value[0].app_id
        rec.app_name = value[0].app_appname
        rec.group_groupname = value[0].group_groupname
        rec.group_description = value[0].group_description
        rec.items = value
        arr.push(rec)
    }
    return arr
}




// G R O U P S   *******************************************************************

function formListGroupSubmit(event) {
    if (event) event.preventDefault()
    model.groups = null
    let search = document.querySelector("#formListGroup input[name='search']").value
    var query =`
    query {
        list_group(
        search: "${search}",
        order: "description ASC"
        ) {
            length
            list {
              description
              groupname
              id
           }
          }
        }    `

    function onSuccess(res){
        model.groups = res.data.list_group.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}



function createGroup(event) {
    if (event) event.preventDefault()
    let groupname =     document.querySelector("#formGroup input[name='groupname']").value
    let description = document.querySelector("#formGroup input[name='description']").value
    
    var query =`
    mutation {
        create_group(
        groupname: "${groupname}",
        description: "${description}"
        ) {
            description
            groupname
            id
          }

        }
    `
    function onSuccess(res){
        alert("create_group success")
        refreshData()
        model.group = res.data["create_group"]
        getGroup(model.group.id)
    } 

    doGraphQLRequest(query, onSuccess, "groupError")
    return false       
}


function updateGroup(event) {
    if (event) event.preventDefault()
    let groupname =     document.querySelector("#formGroup input[name='groupname']"    ).value
    let description = document.querySelector("#formGroup input[name='description']").value
    let id =      document.querySelector("#formGroup input[name='id']"     ).value
   
    var query =`
    mutation {
        update_group(
        id: ${id},
        groupname: "${groupname}",
        description: "${description}"
        ) {
            description
            groupname
            id
          }

        }
    `
    function onSuccess(res){
        alert("update_group success")
        refreshData()
        model.group = res.data["update_group"]
        getGroup(id)
    } 

    doGraphQLRequest(query, onSuccess, "groupError")
    return false       
}



function getGroup(group_id) {
    model.group = null

    var query =`
    query {
        get_group(
        id: ${group_id}
        ) {
            description
            groupname
            id
          }

        list_group_app_role(
        group_id:${group_id}
        ) {
            app_appname
            app_description
            app_id
            app_url
            group_description
            group_groupname
            group_id
            rolename
           }          
                




        list_group_user_role(
        group_id: ${group_id}
        ) {
            group_description
            group_groupname
            group_id
            rolename
            user_description
            user_disabled
            user_email
            user_fullname
            user_id                                
          }
    }          

    `

    function onSuccess(res){
        var group = res.data.get_group
        group.apps = appsOfTheGroup(res.data.list_group_app_role)
        group.users = usersOfTheGroup(res.data.list_group_user_role)
        // render
        model.group = group
    } 

    doGraphQLRequest(query, onSuccess, "groupError")
    return false       
}

// Group array of records by fieldName.
// Returns a map fieldName->array of records
function groupByField(records, fieldName){
    let values = {}
    for (let rec of records ){
        let val = rec[fieldName]
        if (!values[val]) {
            values[val] =[]
        }
        values[val].push(rec)
    }
    return values
}


function appsOfTheGroup(list_group_app_role) {
    let apps = groupByField(list_group_app_role, 'app_id')

    // преобразуем хэш в массив для отображения в mustache
    let arr = []

    for (let [key, value] of Object.entries(apps)) {
        let rec = {}
        rec.app_id =key
        rec.group_id = value[0].group_id
        rec.app_appname = value[0].app_appname
        rec.app_description = value[0].app_description
        rec.app_url = value[0].app_url
        rec.items = value
        arr.push(rec)
    }
    return arr
}

function usersOfTheGroup(list_group_user_role) {
    let users = groupByField(list_group_user_role, 'user_id')

    // преобразуем хэш в массив для отображения в mustache
    let arr = []

    for (let [key, value] of Object.entries(users)) {
        let rec = {}
        rec.user_id =value[0].user_id
        rec.group_id = value[0].group_id
        rec.user_email = value[0].user_email
        rec.user_description = value[0].user_description
        rec.user_fullname = value[0].user_fullname
        rec.items = value
        arr.push(rec)
    }
    return arr
}



function deleteGroup(id) {
    if (!confirm(`Удалить группу ${id}?`)) return false

    var query =`
    mutation {
        delete_group(
        id: ${id}
        ) {
            description
          }
        }
    `
    function onSuccess(res){
        alert("delete_group success")
        model.group = null
        refreshData()
        searchGroups(); 
        showPage('groups') ;
    } 

    doGraphQLRequest(query, onSuccess, "groupError")
    return false       
}




// A P P   U S E R   R O L E   **************************************************************************************************

function getAllApps(event) {
    if (event) event.preventDefault()
    model.allApps = null
    var query =`
    query {
        list_app(
        order: "appname ASC"
        ) {
            length
            list {
              id
              appname
              description
              url
              rebase
              public 
              sign
              xtoken
            }
          }
        }    `

    function onSuccess(res){
        model.allApps = res.data.list_app.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function getAllUsers(event) {
    if (event) event.preventDefault()
    model.allUsers = null
    var query =`
    query {
        list_user(
        order: "fullname ASC"
        ) {
            length
            list {
              username
              fullname
              email
              description
              disabled
              id
              pinhash
              pinrequired
              pinhash_temp    
              emailhash
              emailconfirmed      
              }
          }
        }    `

    function onSuccess(res){
        model.allUsers = res.data.list_user.list
        sortUsersBy('id'         , false, true)
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function getAllGroups(event) {
    if (event) event.preventDefault()
    model.allGroups = null
    var query =`
    query {
        list_group(
        order: "groupname ASC"
        ) {
            length
            list {
              groupname
              description
              id
            }
          }
        }    `

    function onSuccess(res){
        model.allGroups = res.data.list_group.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


function modifyRole(action,appname,username,rolename, onsuccess ) {
    if (!appname || !username || !rolename) return

    var query =`
    mutation {
        ${action}_app_user_role(
        appname: "${appname}",
        username: "${username}",
        rolename: "${rolename}"
        ) {
            rolename
            appname
            username
          }
        }
    `
    function onSuccess(res){
        if (onsuccess) onsuccess()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}

function modifyGroupAppRole(action,group_id,app_id,rolename, onsuccess ) {
    if (!group_id || !app_id || !rolename) 
        return
    var query =`
    mutation {
        ${action}_group_app_role(
        group_id: ${group_id},
        app_id: ${app_id},
        rolename: "${rolename}"
        ) {
            app_appname
            app_description
            app_id
            app_url
            group_description
            group_groupname
            group_id
            rolename
          }
        }
    `
    function onSuccess(res){
        if (onsuccess) onsuccess()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}

function modifyUserGroupRole(action,user_id,group_id,rolename, onsuccess ) {
    if (!user_id || !group_id || !rolename) 
        return
    var query =`
    mutation {
        ${action}_group_user_role(
        user_id: ${user_id},
        group_id: ${group_id},
        rolename: "${rolename}"
        ) {
            group_description
            group_groupname
            group_id
            rolename
            user_description
            user_disabled
            user_email
            user_fullname
            user_id
           }
    }
    `
    function onSuccess(res){
        if (onsuccess) onsuccess()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}

function modifyAppGroupRole(action,app_id,group_id,rolename, onsuccess ) {
    if (!app_id || !group_id || !rolename) 
        return
    var query =`
    mutation {
        ${action}_group_app_role(
        app_id: ${app_id},
        group_id: ${group_id},
        rolename: "${rolename}"
        ) {
            app_appname
            app_description
            app_id
            app_url
            group_description
            group_groupname
            group_id
            rolename
          }
    }
    `
    function onSuccess(res){
        if (onsuccess) onsuccess()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}


function filterRows(selector, value ){
    var v = value.toLowerCase()
    var rows = document.querySelectorAll(selector)
    rows.forEach(e => {
        var txt = e.innerText.toLowerCase()
        if (txt.indexOf(v) == -1) {
            e.classList.add("hidden")
        } else {
            e.classList.remove("hidden")
        }
    });
}

function toggleElements(selector) {
    document.querySelectorAll(selector).forEach(e => e.classList.toggle("hidden"));
}

function hideElements(selector) {
    document.querySelectorAll(selector).forEach(e => e.classList.add("hidden"));
}


function showElements(selector) {
    document.querySelectorAll(selector).forEach(e => e.classList.remove("hidden"));
}

function addClass(selector, classname) {
        document.querySelectorAll(selector).forEach(e => e.classList.add(classname));
}

function removeClass(selector, classname) {
    document.querySelectorAll(selector).forEach(e => e.classList.remove(classname));
}


function hidePassword() {
    document.querySelector("#formUser input[name='password']").setAttribute('type','password')
}

function showPassword() {
    document.querySelector("#formUser input[name='password']").setAttribute('type','text')
}

function removeQueryStringFromBrowserURL() {
    let url = document.location.href
    let url1 = url.replace(/\?.*/, '')
    let url2 = url1+"?url="+model.appurl
    history.replaceState(null,null,url2)   
}

function newPassword (n, numbers=false) {
    let pickSymbol =(s) => s[Math.floor(Math.random()*s.length)]
    var symbolSets =numbers? ["0123456789"]:["bcdfghjklmnpqrstvwxz","aeiou"] 
    var password = ''
    for (let i=0; i<n; i++){
        password += pickSymbol(symbolSets[i%symbolSets.length])
    }
    return password
}


function generatePassword() {
    document.querySelector("#formUser input[name='password']").value = newPassword(9)
}

function generatePinHash() {
    let hash = (window.crypto && window.crypto.randomUUID) ? window.crypto.randomUUID() : newPassword(20,true);
    document.querySelector("#formUser input[name='pinhash_temp']").value = hash;
    setPinUrl();
}

function setPinUrl(){
    let input = document.querySelector("#formUser input[name='pinhash_temp']")
    let hash = input.value
    hash = hash.replaceAll(" ","")
    input.value = hash
    // let url = `/set-authenticator.html#username=${model.user.username}&hash=${model.user.pinhash_temp}&authurl=${model.appurl}`
    // let link = document.getElementById('newPinUrl')
    // link.innerText = url
    // link.href = url
    document.getElementById('pinUrlContainer').style.display= "none"
}

function buildSocialIcons(url) {
        fetch(url).then(x => x.json())
        .then( t => renderOauthProvidersJSON(t) )
        .catch( err => {
            console.log("fetch error:",err)
            return
            }
        )  
}

function renderOauthProvidersJSON(jsn) {
    let el = document.getElementById('socialIcons')
    
    if (!el) return
    
    el.innerHTML = ''
    let icons = []
    for (let [k,v] of Object.entries(jsn) ) {
        icons.push(`
        <div class="social-element">
        <a class="button button-clear" title="login via ${k}" href="${model.appurl + v[0]}">
        ${k}
        <br>
        <img class="social-icon" src="images/${k}.svg">
        </a>
        <!--
        <br>
        <a class="social-logout" title="logout from ${k}" href="${model.appurl + v[1]}">logout</a>
        -->
        </div>
        
        
        `)
    }
    
    
    if (icons.length > 0){
        el.innerHTML = '<div class="socialHeader">войти через социальную сеть</div>' + icons.join(' ')
        showElements("#selfRegHelpText")
        removeClass('#socialIcons', 'transparent')
    } else {
        hideElements("#selfRegHelpText")
    }
}



function getNewCaptcha() {
    let uri = model.appurl+"/captcha?"+ new Date().getTime()
    document.getElementById("captchaImg").src = uri
    return false
}


// G O O G L E   C H A R T S  ***************************************************************

// google.charts.load('current', {packages: ['gauge']})
// google.charts.setOnLoadCallback(drawGauge)


var gauges = {}

function drawGauge(title, val, maxV, containerID) {

    if (!google) return
    if (!google.visualization) return
    if (!google.visualization.arrayToDataTable) return
    if (!google.visualization.Gauge) return
    
    var container = document.getElementById(containerID)
    if (! container) return

    if (! gauges[title]) {
        gauges[title]={}
        gauges[title].data = google.visualization.arrayToDataTable([['Label', 'Value'], [title, 0]])
    }
    if (container.innerHTML == ""){
        gauges[container]=container
        gauges[title].chart = new google.visualization.Gauge(container)
        console.log("Redraw gauge container:", containerID)
    }

    // set values
    gauges[title].data.setValue(0, 1, val)

    let maxVal = maxV ? maxV: getUpperLimit(val)
    var options = {
        // width: 120,
        height:  120, 
        animation:{
            duration: 700
        },
        greenFrom:0, greenTo: maxVal*0.2,
        yellowFrom: maxVal*0.8, yellowTo: maxVal*0.9,
        redFrom: maxVal*0.9, redTo: maxVal,
        minorTicks: 5,
        // redColor: '#E10098',
        // greenColor: 'cyan',
        max: maxVal
    }

    // draw the chart
    gauges[title].chart.draw(gauges[title].data, options)
}


function getUpperLimit(n) {
    var lim = 10
    while (n > lim) lim *= 10
    return lim
}

// **********************************************************************************************
// **********************************************************************************************
// **********************************************************************************************


function clearLoginForm() {
    document.getElementById("loginUsername").value = ""
    document.getElementById("loginPassword").value = ""
    document.getElementById("loginCaptcha").value = ""
    document.getElementById("loginPin").value = ""
    document.getElementById("loginError").innerText = ""
    document.getElementById("socialLoginError").innerHTML = ""
    document.getElementById("oauth2email").innerHTML = ""
}


function logout() {
    logoutGraphQLFormSubmit()
    clearLoginForm()
    showPage('login',true)
    isSelfRegAllowed()
    model.captchaRequired = false
    return false
}


// refreshData() если пользователь залогинен наполняет модель данными,
// в противном случае или обнуляет модель.
function refreshData() {
    if (model.logined) {
        getAllApps()
        getAllUsers()
        getAllGroups()
        formListAppSubmit()
        formListUserSubmit()  
        formListGroupSubmit()  
    } else {
        // nullify model's inner props
        for (const k of Object.keys(model)) {
            if (k.startsWith('_')) {
                model[k] = null
            }
        }
        isSelfRegAllowed()
   }    
}

function getCurrentPageID(){
    var p = location.hash.slice(1)
    return (!p || p.includes('=')) ? 'apps' : p
}

// refreshApp обновляет данные модели и GUI
function refreshApp() {
    
    refreshData()
    
    if (model.logined) {
        let page = getCurrentPageID()
        showPage(page) 
        showElements('#menu')     
        hideElements('#menuUnlogined')     
    } else {
        showPage('login',true)
        hideElements('#menu')
        showElements('#menuUnlogined')
    }  
    // предзагружаем шаблоны страниц
    renderTemplateFile('mustache/group.html', model, '#groupPage')
    renderTemplateFile('mustache/app.html', model, '#appPage')
    renderTemplateFile('mustache/user.html', model, '#userPage')
}



window.onhashchange = function(event) {
    console.debug("onhashchange ", event)
    var newpage = event.newURL.split('#')[1]
    if (newpage) 
        showPage(newpage)
}


function setAppParams(){
    model.oauth2email = model.urlParams.get('oauth2email') || ''
    model.oauth2name = model.urlParams.get('oauth2name') || ''
    model.oauth2id = model.oauth2email.replace(/@.*/,'')
    document.getElementById('oauth2email').innerHTML = model.oauth2email

    var url = model.urlParams.get('url')
    var css = model.urlParams.get('css')
    var oauth2error = model.urlParams.get('oauth2error')
    var alrt = model.urlParams.get('alert')
    if (alrt) alert(alrt)
    document.getElementById('socialLoginError').innerHTML = oauth2error
    if (css) 
        document.getElementById('theme').href = css
    model.appurl = url ? url : 'https://auth-proxy.rg.ru'
    removeQueryStringFromBrowserURL()
    //?
    model.captchaRequired = false
}



function init() {
    renderTemplateFile('mustache/params.html', model, '#paramsPage')
    if (google && google.charts) {
        google.charts.load('current', {'packages':['gauge']})
        google.charts.setOnLoadCallback(getAppstatRest)
    }
    setAppParams()
    getLoginedUser()
    refreshApp()   
    getParams()   
}


// O N   P A G E   L O A D  ****************************************************************************************

init()
