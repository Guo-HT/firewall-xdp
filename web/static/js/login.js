$(function () {
    var sys_info = get_system_banner()
    $("#proj_title").text(sys_info.title)

    layui.use(['form', "layer"], function () {
        var form = layui.form;
        var layer = layui.layer;

        //监听提交
        $("#login-btn").click(function (data) {
            $.ajax({
                url: "/user/status/login",
                type: "post",
                data: JSON.stringify({
                    username:$("#username-input").val(),
                    password:$("#password-input").val(),
                }),
                dataType: "json"
            }).done(function(msg){
                if (msg.code===200){
                    if (msg.data.isDefaultPwdChanged!==true){
                        layer.msg("首次登录需修改密码")

                        layer.open({
                            type:1,
                            area: ["400px", "350px"],
                            title: "首次登录需修改密码",
                            content: $("#change_password_form"),
                            shade: 0,
                            btn:["确认", "重置"],
                            btn1: function(index, layero){
                                // console.log("提交")
                                var password_old = $("#change-password-old").val();
                                var password_new = $("#change-password-input").val();
                                var password_new_check = $("#change-password-check-input").val();
                                // console.log(password_old, password_new, password_new_check)
                                if (password_new!==password_new_check){
                                    layer.msg("密码校验错误")
                                    return
                                }
                                if(password_old===password_new){
                                    layer.msg("新旧密码不能一致")
                                    return
                                }
                                $.ajax({
                                    url:"/user/status/changePwd",
                                    type:"post",
                                    data:JSON.stringify({
                                        username: $("#username-input").val(),
                                        old_password:password_old,
                                        new_password:password_new,
                                    }),
                                    dataType:"json"
                                }).done(function(msg){
                                    if (msg.code===200){
                                        layer.msg("修改成功，请重新登陆")
                                        $("#password-input").val("")
                                        layer.closeAll();
                                    }else{
                                        layer.msg(msg.msg)
                                    }
                                }).fail(function(e){
                                    layer.msg("error")
                                })
                            },
                            btn2: function(index, layero){
                                // console.log("重置")
                                $("#change-password-reset").click()
                                return false
                            },
                            cancel: function(layero, index){
                                // console.log("cancel")
                                layer.closeAll();
                            }
                        })
                    }else{
                        window.location.href="/"
                        // layer.msg("不需要修改密码")
                    }
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
        });



    });
})