$(function () {
    layui.use(['element', 'table', "form", "layer", "slider"], function () {
        const element = layui.element;
        const table = layui.table;
        const form = layui.form;
        const layer = layui.layer;
        const slider = layui.slider;

        let slider_start;
        let slider_end;

        render_protocol_table()
        get_protocol_engine_status()

        // 传入网卡名，将IP黑名单渲染到表格
        function render_protocol_table() {
            table.render({
                elem: "#protocol_table",
                url: "/func/proto/rules",
                method: "get",
                page: false,
                parseData: function (res) { //res 即为原始返回的数据
                    var interface_data = []
                    if (res.code === 200) {
                        if (res.data === null || res.data.length === 0) {
                            return {
                                "code": res.code === 200 ? 0 : -1, //解析接口状态
                                "msg": res.msg, //解析提示文本
                                "count": res.data.total, //解析数据长度
                                "data": [{
                                    "id": "-",
                                    "protocol_name": "-",
                                    "req_type": "-",
                                    "req_regx": "-",
                                    "rsp_type": "-",
                                    "rsp_regx": "-",
                                    // "start_port": "-",
                                    // "end_port": "-",
                                    "port_range": "-"
                                    // "is_enable": "-",
                                }] //解析数据列表
                            };
                        }
                        for (var i = 0; i < res.data.length; i++) {
                            var is_enable_content = "    <input type=\"checkbox\" name=\"protocol_switch\" lay-skin=\"switch\" lay-event='toggle_switch' lay-text=\"开启|关闭\" "
                            if (res.data[i].is_enable === true) {
                                is_enable_content += "checked>"
                            } else {
                                is_enable_content += ">"
                            }
                            interface_data.push({
                                "id": res.data[i].id,
                                "protocol_name": res.data[i].protocol_name,
                                "req_type": res.data[i].req_type,
                                "req_regx": res.data[i].req_regx,
                                "rsp_type": res.data[i].rsp_type,
                                "rsp_regx": res.data[i].rsp_regx,
                                // "start_port": res.data[i].start_port,
                                // "end_port": res.data[i].end_port,
                                "port_range": res.data[i].start_port + " ~ " + res.data[i].end_port,
                                "is_enable": res.data[i].is_enable,
                            })
                        }
                        return {
                            "code": res.code === 200 ? 0 : -1, //解析接口状态
                            "msg": res.msg, //解析提示文本
                            "count": res.data.length, //解析数据长度
                            "data": interface_data //解析数据列表
                        };
                    } else {
                        return {
                            "code": res.code === 200 ? 0 : -1, //解析接口状态
                            "msg": res.msg, //解析提示文本
                            "count": res.data.total, //解析数据长度
                            "data": [] //解析数据列表
                        };
                    }
                },
                cols: [[
                    {field: 'id', title: 'ID', width: "4%", align: "center", sort: false},
                    {field: 'protocol_name', title: '协议名', width: "10%", align: "center", sort: false},
                    {field: 'req_type', title: "请求类型", width: "10%", sort: false, align: "center", toolbar: ""},
                    {field: 'req_regx', title: "请求特征", width: "15%", sort: false, align: "center", toolbar: ""},
                    {field: 'rsp_type', title: "响应类型", width: "10%", sort: false, align: "center", toolbar: ""},
                    {field: 'rsp_regx', title: "响应特征", width: "15%", sort: false, align: "center", toolbar: ""},
                    // {field: 'start_port', title: "端口起始", width: "8%", sort: false, align: "center", toolbar:""},
                    // {field: 'end_port', title: "端口中止", width: "8%", sort: false, align: "center", toolbar:""},
                    {field: 'port_range', title: "端口范围", width: "16%", sort: false, align: "center", toolbar: ""},
                    {
                        field: 'is_enable',
                        title: "启用状态",
                        width: "10%",
                        sort: false,
                        align: "center",
                        templet: "#protocol_switch"
                    },
                    {
                        field: 'option',
                        title: "操作",
                        width: "10%",
                        sort: false,
                        align: "center",
                        toolbar: "#protocol_opt"
                    },
                ]],
            })
        }

        // 获取协议分析引擎开关状态
        function get_protocol_engine_status(){
            $.ajax({
                url: "/func/proto/status",
                type: "get",
                dataType: "json",
            }).done(function(msg){
                if(msg.code===200){
                    var protocol_status = msg.data
                    if (protocol_status){
                        // console.log("开启")
                        $('#protocol_engine_switch').attr({"checked": 'checked'});
                    }else{
                        // console.log("关闭")
                        $('#protocol_engine_switch').removeAttr("checked").attr("value", "off")
                    }
                    form.render("checkbox")
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
        }

        // 监听删除事件
        table.on("tool(protocol_table)", function (obj) {
            var del_rule_id = obj.data.id
            $.ajax({
                url: "/func/proto/delRules",
                type: "post",
                data: JSON.stringify({
                    id: del_rule_id
                }),
                dataType: "json",
                contentType: "application/json;charset=utf-8",
                async: false
            }).done(function(msg){
                if(msg.code===200){
                    render_protocol_table()
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
        })

        // 监听某协议开关事件
        form.on("switch(switch)", function (obj) {
            // console.log(obj)
            var status;
            var belong_protocol = $(this).attr("belong");
            if (this.value === "true") {
                this.value = "false"
                status = false
            } else {
                this.value = "true"
                status = true
            }
            // console.log(this.value)
            $.ajax({
                url: "/func/proto/rules",
                type: "post",
                data: JSON.stringify({
                    "protoName": belong_protocol,
                    "status": status
                }),
                dataType: "json",
                contentType: "application/json;charset=utf-8",
                async: false
            }).done(function (msg) {
                if (msg.code !== 200) {
                    layer.msg(msg.msg)
                }
                // render_protocol_table()
            }).fail(function (msg) {
                layer.msg("error")
            })
        })

        // 监听协议分析引擎开关
        form.on("switch(protocol_engine_switch)", function(obj){
            var target_status = $(this).is(':checked')
            var url = target_status? "/func/proto/start": "/func/proto/stop"

            // console.log(target_status, url)
            $.ajax({
                url: url,
                type: "post",
                dataType:"json",
            }).done(function(msg){
                if (msg.code === 200){
                    form.render("checkbox")
                }else{
                    layer.msg("协议分析功能开关切换失败")
                }
            }).fail(function(e){
                layer.msg("error")
            })

        })

        // 监听添加事件
        $("#add-protocol").click(function () {
            var slider_port_range = slider.render({
                elem: "#port-range-slider",
                min: 0,
                max: 65535,
                range: true,
                value: [0, 65535],
                change: function (value) {
                    slider_start = value[0]
                    slider_end = value[1]
                }
            })
            layer.open({
                type: 1,
                area: ['66%', '400px'],
                title: '请输入一条策略',
                content: $("#input-feature-from"),
                shade: 0,
                btn: ['提交', '重置']
                , btn1: function (index, layero) {
                    // console.log("提交")
                    var protocol_name = $("#protocol_name").val()
                    var req_type = $("#req_type input:checked").val()
                    var req_reg = $("#req_reg").val()
                    var rsp_type = $("#rsp_type input:checked").val()
                    var rsp_reg = $("#rsp_reg").val()
                    if (
                        protocol_name === "" ||
                        typeof (req_type) === undefined ||
                        req_reg === "" ||
                        typeof (rsp_type) === undefined ||
                        rsp_reg === ""
                    ) {
                        layer.msg("请完善信息")
                    } else {
                        // console.log("ok")
                        if (slider_start === undefined || slider_end === undefined) {
                            slider_start = 0
                            slider_end = 65535
                        }
                        var new_rule = {
                            "id": 9999999,
                            "protocol_name": protocol_name,
                            "req_type": req_type,
                            "req_regx": req_reg,
                            "rsp_type": rsp_type,
                            "rsp_regx": rsp_reg,
                            "start_port": slider_start,
                            "end_port": slider_end,
                            "is_enable": false
                        }
                        // console.log(new_rule)
                        $("#protocol_rule_reset").click()
                        $.ajax({
                            url: "/func/proto/addRules",
                            type: "post",
                            data: JSON.stringify(new_rule),
                            dataType: "json",
                            contentType: "application/json;charset=utf-8",
                            async: false
                        }).done(function(msg){
                            if(msg.code===200){
                                $("#protocol_rule_reset").click()
                                render_protocol_table()
                            }else{
                                layer.msg(msg.msg)
                            }
                        }).fail(function(e){
                            layer.msg("error")
                        })

                        layer.closeAll();
                    }


                },
                btn2: function (index, layero) {
                    // console.log("重置")
                    $("#protocol_rule_reset").click()
                    return false
                    // layer.closeAll();
                },
                cancel: function (layero, index) {
                    // console.log("cancel")
                    layer.closeAll();
                }
            });


        })

        // 监听重载协议规则事件
        $("#reload-protocol").click(function(){
            $.ajax({
                url: "/func/proto/reloadRules",
                type: "post",
                dataType: "json"
            }).done(function(msg){
                if (msg.code===200){
                    render_protocol_table()
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function(e){
                layer.msg("error")
            })
        })
    })
})