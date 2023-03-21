function get_netcard_engine_list() {
    var netcard_list = []
    $.ajax({
        url: "/iface/engine/getEngineList",
        type: "get",
        dataType: "json",
        async: false,
    }).done(function (msg) {
        if (msg.code === 200) {
            if (msg.data != null) {
                netcard_list = msg.data
            } else {
                netcard_list = []
            }
        } else {
            netcard_list = null
        }
    })
    return netcard_list
}

function render_netcard_list() {
    var netcard_list = get_netcard_engine_list()
    // console.log(netcard_list)
    if (netcard_list !== null) {
        if (netcard_list.length!==0) {
            for (var i = 0; i < netcard_list.length; i++) {
                $("#netcard-select-input").append("<option value='" + netcard_list[i] + "'>" + netcard_list[i] + "</option>")
            }
        }else{
            $("#netcard-select-input").append("<option value='default'>请选择一张网卡</option>")
        }
    } else {
        $("#netcard-select-input").append("<option value='default'>请选择一张网卡</option>")
    }
    layui.use("form", function () {
        var form = layui.form;
        form.render("select");
    })
}

function validateIP(ip) {
    var regexIP = /^((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))$/g;
    if (regexIP.test(ip)) {
        return true;
    } else if(ip!==''){
        //alert("ip地址有误，请重新输入！");
        return false;
    }
}


function get_system_banner(){
    var title = "";
    var icon = "";
    $.ajax({
        url:"/status/setting/systemTitle",
        type:"get",
        dataType:"json",
        async: false,
    }).done(function(msg){
        if (msg.code===200){
            title = msg.data.title
            icon = msg.data.icon
        }
    }).fail(function(e){
        title = null
        icon = null
    })
    return {"title": title, "icon": icon}
}