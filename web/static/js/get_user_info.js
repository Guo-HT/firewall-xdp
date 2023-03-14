get_user_login_status()

function get_user_login_status() {
    $.ajax({
        url: "/user/status/info",
        type: "get",
        dataType: "json",
    }).done(function (msg) {
        // console.log(msg)
        if (msg.code === 200) {
            if (msg.data.loginState === false) {
                window.parent.location.href = "/login";
                window.location.href = "/login";
            }
        }
    }).fail(function (e) {
        ;
    })
}
