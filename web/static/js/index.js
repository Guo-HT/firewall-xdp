window.onload = function () {
    var sys_info = get_system_banner()
    $("#proj_title").text(sys_info.title)

    layui.use(['element', 'jquery', "layer"], function () {
        var element = layui.element;
        var layer = layui.layer;


        $(".my-page-nav").click(function () {
            var target = $(this).attr("goto")
            goto_page(target)
        })

        // 刷新页面时，获取最后一级，直接跳转
        goto_page(get_last_url())

        // 跳转
        function goto_page(target) {
            history.pushState({}, null, location.origin + location.pathname + "#/" + target);
            $("#my-main-iframe").attr("src", "/web/static/page/" + target + ".html")
        }

        // 获取最后一级uri
        function get_last_url() {
            var url = window.location.href;
            var index = url.lastIndexOf("/");
            var str = url.substring(index + 1)
            return str === "" ? "overview" : str
        }

        $("#logout_btn").click(function () {
            layer.confirm("确认注销？", {
                btn: ['确认', "取消"],
                btn1: function (index, layero) {
                    $.ajax({
                        url: "/user/status/logout",
                        type: "post",
                        dataType: "json"
                    }).done(function(msg){
                        if (msg.code===200){
                            window.location.href = "/login"
                        }else{
                            layer.msg(msg.msg)
                        }
                    }).fail(function(e){
                        layer.msg("error")
                    })
                    layer.closeAll()
                },
                btn2: function(index, layero){
                    layer.closeAll()
                }
            })
        })

    });
}


