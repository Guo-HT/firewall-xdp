$(function () {
    var sys_info = get_system_banner()

    layui.use(['element', 'table', "form", "layer", "slider", "upload"], function () {
        const element = layui.element;
        const table = layui.table;
        const form = layui.form;
        const layer = layui.layer;
        const slider = layui.slider;
        const upload = layui.upload;

        $("#sys-title").val(sys_info.title)

        render_user_table()

        function render_user_table() {
            table.render({
                elem: "#user_table",
                url: "/user/status/allUser",
                method: 'get',
                page: {
                    limit: 10,
                    limits: [10, 20, 30],
                },
                request: {
                    pageName: 'page_no' //页码的参数名称，默认：page
                    , limitName: 'page_size' //每页数据量的参数名，默认：limit
                },
                parseData: function (res) {
                    var interface_data = []
                    if (res.code === 200) {
                        if (res.data === null || res.data.length === 0) {
                            return {
                                "code": res.code === 200 ? 0 : -1,
                                "msg": res.msg,
                                "count": res.data.data.total,
                                "data": [{
                                    "id": "-",
                                    "username": "-",
                                    "email": "-",
                                    "role": "-",
                                    "create_at": "-",
                                }]
                            };
                        }
                        for (var i = 0; i < res.data.data.length; i++) {
                            var role = res.data.data[i].role;
                            var roleStr = ""
                            if (role === 0) {
                                roleStr = "管理员"
                            } else if (role === 1) {
                                roleStr = "操作员"
                            } else {
                                roleStr = "审计员"
                            }

                            interface_data.push({
                                "id": res.data.data[i].id,
                                "username": res.data.data[i].username,
                                "email": res.data.data[i].email,
                                "role": roleStr,
                                "create_at": res.data.data[i].create_at,
                            })
                        }
                        return {
                            "code": res.code === 200 ? 0 : -1, //解析接口状态
                            "msg": res.msg, //解析提示文本
                            "count": res.data.data.length, //解析数据长度
                            "data": interface_data //解析数据列表
                        }
                    }
                }
                , cols: [[
                    {field: 'id', title: 'ID', width: "10%", sort: false, align: "center"},
                    {field: 'username', title: '用户名', width: "20%", sort: false, align: "center"},
                    {field: 'email', title: '邮箱', width: "20%", sort: false, align: "center"},
                    {field: 'role', title: '角色', width: "10%", sort: false, align: "center"},
                    {field: 'create_at', title: '创建时间', width: "20%", sort: false, align: "center"},
                    {field: 'option', title: '操作', width: "20%", toolbar: "#user_opt", sort: false, align: "center"},
                ]]
            })
        }

        // 监听表格操作事件
        table.on("tool(user_table)", function (obj) {
            // console.log(obj)
            var this_data = obj.data;
            var event = obj.event;
            var username = this_data.username
            if (event === "change_pwd") {
                layer.open({
                    type: 1,
                    area: ["400px", "350px"],
                    title: "修改密码",
                    content: $("#change_password_form"),
                    shade: 0,
                    btn: ["确认", "重置"],
                    btn1: function (index, layero) {
                        // console.log("提交")
                        var password_old = $("#change-password-old").val();
                        var password_new = $("#change-password-input").val();
                        var password_new_check = $("#change-password-check-input").val();
                        // console.log(password_old, password_new, password_new_check)
                        if (password_new !== password_new_check) {
                            layer.msg("密码校验错误")
                            return
                        }
                        if (password_old === password_new) {
                            layer.msg("新旧密码不能一致")
                            return
                        }
                        $.ajax({
                            url: "/user/status/changePwd",
                            type: "post",
                            data: JSON.stringify({
                                username: username,
                                old_password: password_old,
                                new_password: password_new,
                            }),
                            dataType: "json"
                        }).done(function (msg) {
                            if (msg.code === 200) {
                                layer.msg("修改成功，请重新登陆")
                                // $("#password-input").val("")
                                window.parent.location.href = "/login";
                                layer.closeAll();
                            } else {
                                layer.msg(msg.msg)
                            }
                        }).fail(function (e) {
                            layer.msg("error")
                        })
                    },
                    btn2: function (index, layero) {
                        // console.log("重置")
                        $("#change-password-reset").click()
                        return false
                    },
                    cancel: function (layero, index) {
                        // console.log("cancel")
                        layer.closeAll();
                    }
                })
            } else if (event === "delete") {
                // alert("删除用户还没做");

                ////////////////
                layer.open({
                    type: 1,
                    area: ["400px", "200px"],
                    title: "删除",
                    content: $("#del_user_form"),
                    shade: 0,
                    btn: ["确认", "重置"],
                    btn1: function (index, layero) {
                        // console.log("提交")
                        var password = $("#del_user_input").val();
                        $.ajax({
                            url: "/user/status/delUser",
                            type: "post",
                            data: JSON.stringify({
                                "password": password,
                                "target_user_name": username,
                            })
                        }).done(function (msg) {
                            if (msg.code === 200) {
                                render_user_table();
                                layer.closeAll();
                            } else {
                                layer.msg(msg.msg);
                            }
                        }).fail(function (e) {
                            layer.msg("error")
                        })

                    },
                    btn2: function (index, layero) {
                        // console.log("重置")
                        $("#change-password-reset").click()
                        return false
                    },
                    cancel: function (layero, index) {
                        // console.log("cancel")
                        layer.closeAll();
                    }
                })
            }
        })

        // 监听添加用户事件
        $("#add-user").click(function () {
            // alert();
            layer.open({
                type: 1,
                area: ['400px', '480px'],
                title: '请输入用户信息',
                content: $("#new-user-from"),
                shade: 0,
                btn: ['提交', '重置']
                , btn1: function (index, layero) {
                    var cur_password = $("#cur-pwd").val();
                    var new_username = $("#new-username").val();
                    var new_pwd = $("#new-pwd").val();
                    var new_pwd_check = $("#new-pwd-check").val();
                    var new_email = $("#new-email").val();
                    var new_role = parseInt($("#new-role input[type='radio']:checked").val());
                    console.log(cur_password, new_username, new_pwd, new_pwd_check, new_email, new_role)
                    if (new_pwd !== new_pwd_check) {
                        layer.msg("校验密码不一致")
                        return
                    }
                    $.ajax({
                        url: "/user/status/addUser",
                        type: "post",
                        data: JSON.stringify({
                            "cur_user_password": cur_password,
                            "user_name": new_username,
                            "password": new_pwd,
                            "email": new_email,
                            "role": new_role,
                        }),
                        dataType: "json"
                    }).done(function (msg) {
                        if (msg.code === 200) {
                            render_user_table();
                            layer.closeAll();
                        } else {
                            layer.msg(msg.msg)
                        }
                    }).fail(function (e) {
                        layer.msg("error")
                    })
                }
                , btn2: function (index, layero) {
                    $("#add-user-reset").click()
                    return false
                }
                , cancel: function (layero, index) {
                    layer.closeAll();
                }
            })
        })

        render_network_table()

        // 渲染网卡状态信息
        function render_network_table() {
            table.render({
                elem: "#netcard_table",
                url: "/iface/engine/getIfaceList",
                method: 'get',
                page: false,
                parseData: function (res) {
                    var interface_data = []
                    if (res.code === 200) {
                        if (res.data === null || res.data.length === 0) {
                            return {
                                "code": res.code === 200 ? 0 : -1,
                                "msg": res.msg,
                                "count": 1,
                                "data": [{
                                    "netcard_name": "-",
                                    'ip': "-",
                                    "mac": "-",
                                    "flags": "-",
                                    "is_up": "-",
                                    "is_attach": "-",
                                }]
                            };
                        }
                        for (var i = 0; i < res.data.length; i++) {
                            var ip_html = "";
                            var ip_join = ""
                            if (res.data[i].ip !== null ) {
                                for (var j = 0; j < res.data[i].ip.length; j++) {
                                    ip_html += ('<span class="layui-badge layui-bg-gray">' + res.data[i].ip[j] + '</span>&nbsp;')
                                    ip_join += (res.data[i].ip[j]+"|")
                                }
                            }else{
                                ip_html = "-"
                            }
                            var flag = res.data[i].flags.split("|")
                            var flag_html = ""
                            var is_up  = false;
                            if (flag !== null ) {
                                for (var k = 0; k < flag.length; k++) {
                                    if (flag[k]==="up"){
                                        is_up = true
                                    }
                                    flag_html += ('<span class="layui-badge layui-bg-gray">' + flag[k] + '</span>&nbsp;')
                                }
                            }else{
                                flag_html = "-"
                            }
                            interface_data.push({
                                "netcard_name": res.data[i].netcard_name,
                                'ip': ip_html,
                                "mac": '<span class="layui-badge layui-bg-gray">' + res.data[i].mac + '</span>',
                                "flags": flag_html,
                                "is_up": is_up,
                                "is_attach": res.data[i].engine_status,
                                "ip_join": ip_join,
                            })
                        }
                        return {
                            "code": res.code === 200 ? 0 : -1,
                            "msg": res.msg,
                            "count": res.data.length,
                            "data": interface_data,
                        }
                    }
                },
                cols: [[
                    {field: 'netcard_name', title: '网卡名', width: "10%", sort: false, align: "center"},
                    {field: 'ip', title: 'IP地址', width: "20%", sort: false, align: "center"},
                    {field: 'mac', title: 'MAC地址', width: "20%", sort: false, align: "center"},
                    {field: 'flags', title: '状态', width: "30%", sort: false, align: "center"},
                    {field: 'engine_opt', title: '挂载引擎', width: "10%", toolbar: "#engine_table", sort: false, align: "center"},
                    {field: 'netcard_opt', title: '启停网卡', width: "10%", toolbar: "#netcard_switch", sort: false, align: "center"},
                ]]
            })
        }

        // 开/关网卡
        form.on("switch(netcard_switch_evt)", function(obj){
            var attr_ip = $(this).attr("ip")
            if (attr_ip.includes(window.location.hostname)){
                layer.msg("当前网卡正在为您提供网络服务，无法关闭")
                render_network_table()
                return
            }
            var load_index = layer.load();
            var netcard_name = $(this).attr("belong")
            var engine_status = $(this).attr("engine")
            var target_status = obj.value
            console.log("网卡启停", obj, netcard_name, engine_status)
            if(engine_status==="true" || engine_status===true){
                layer.msg("引擎正在运行，请写在引擎后关闭网卡！")
                render_network_table()
                return
            }
            var url = ""
            if (target_status==="false" || target_status===false){
                url = "/iface/status/startIface"
            }else{

                url = "/iface/status/stopIface"
            }

            $.ajax({
                url: url,
                type: "post",
                data: JSON.stringify({
                    "iface": netcard_name,
                }),
                dataType: "json",
                async: false,
            }).done(function(msg){
                if (msg.code===200){

                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
            // form.render("switch")
            render_network_table()
            layer.close(load_index);
        })

        // 挂载/卸载引擎
        form.on("switch(engine_switch_evt)", function(obj){
            var load_index = layer.load();
            console.log("切换引擎", obj, $(this).attr("belong"))
            var netcard_name = $(this).attr("belong")
            var target_status = obj.value

            var url = ""
            if (target_status==="false" || target_status===false){
                url = "/iface/engine/attach"
            }else{
                url = "/iface/engine/detach"
            }
            $.ajax({
                url: url,
                type: "post",
                data: JSON.stringify({
                    "iface": netcard_name,
                }),
                dataType: "json",
                async: false,
            }).done(function(msg){
                if (msg.code===200){

                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
            // form.render("switch")
            render_network_table()
            layer.close(load_index);
        })

        // 提交系统名称
        $("#submit_sys_info").click(function(){
            var sys_title = $("#sys-title").val()
            $.ajax({
                url:"/status/setting/systemTitle",
                type:"post",
                dataType:"json",
                data: JSON.stringify({
                    "title": sys_title,
                    "icon": "firewall.png"
                })
            }).done(function(msg){
                if (msg.code===200){
                    window.top.location.reload()
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
        })

    })
})