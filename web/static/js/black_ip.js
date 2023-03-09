$(function () {
    layui.use(['element', 'table', "form", "layer"], function () {
        const element = layui.element;
        const table = layui.table;
        const form = layui.form;
        const layer = layui.layer;

        render_netcard_list()  // 渲染网卡选择列表

        var netcard = $("#netcard-select-input").val()
        render_black_ip_table(netcard)

        // 传入网卡名，将IP黑名单渲染到表格
        function render_black_ip_table(netcard) {
            if (netcard === "default") {
                layer.msg("请选择一张网卡")
                // 渲染空表格
                table.render({
                    elem: "#black_ip_table",
                    page: {
                        limit: 10,
                        limits: [10, 20, 30],
                    },
                    cols: [[
                        {field: 'black_ip', title: 'IP黑名单', width: "40%", sort: false},
                        {field: 'hit', title: '命中次数', width: "40%", sort: false},
                        {field: 'option', title: "操作", width: "20%", sort: false, align: "center"},
                    ]],
                    data: [
                        {"black_ip": "-", "hit": '-', "option": "-"}
                    ]
                })
                return
            }
            table.render({
                elem: "#black_ip_table",
                url: "/xdp/black/getIP",
                method: "get",
                request: {
                    pageName: 'page_no' //页码的参数名称，默认：page
                    , limitName: 'page_size' //每页数据量的参数名，默认：limit
                },
                where: {
                    "iface": netcard
                },
                page: {
                    limit: 10,
                    limits: [10, 20, 30],
                },
                parseData: function (res) { //res 即为原始返回的数据
                    var interface_data = []
                    console.log(res.data.data)
                    if (res.code === 200) {
                        if (res.data.data === null) {
                            return {
                                "code": res.code === 200 ? 0 : -1, //解析接口状态
                                "msg": res.msg, //解析提示文本
                                "count": res.data.total, //解析数据长度
                                "data": [{"black_ip": "-", "hit": "-", "option": "-"}] //解析数据列表
                            };
                        }
                        for (var i = 0; i < res.data.data.length; i++) {
                            interface_data.push({
                                "black_ip": res.data.data[i].ip,
                                "hit": res.data.data[i].hit,
                                // "option": '<a class="layui-btn layui-btn-xs delete_ip_black" onclick="func_del_black_ip(\'' + res.data.data[i].ip + '\');">删除</a>',
                            })
                        }
                        return {
                            "code": res.code === 200 ? 0 : -1, //解析接口状态
                            "msg": res.msg, //解析提示文本
                            "count": res.data.total, //解析数据长度
                            "data": interface_data //解析数据列表
                        };
                    } else {
                        return {
                            "code": res.code === 200 ? 0 : -1, //解析接口状态
                            "msg": res.msg, //解析提示文本
                            "count": res.data.total, //解析数据长度
                            "data": "" //解析数据列表
                        };
                    }

                },
                cols: [[
                    {field: 'black_ip', title: 'IP黑名单', width: "40%", sort: false},
                    {field: 'hit', title: '命中次数', width: "40%", sort: false},
                    {field: 'option', title: "操作", width: "20%", sort: false, align: "center", toolbar:"#del_btn"},
                ]],
            })
        }

        // 监听select表单变化事件
        form.on("select(select_changed)", function (data) {
            console.log(data.value)
            render_black_ip_table(data.value)
        })

        // 监听删除按钮
        table.on('tool(black_ip_table)', function(obj){
            // console.log(obj)
            var netcard = $("#netcard-select-input").val()
            var target = [obj.data.black_ip.replace("/32", ""), ]
            var del_data = {
                "iface": netcard,
                "blackIpList": target
            }
            console.log("删除一个 [黑名单IP]", del_data)
            $.ajax({
                url: "/xdp/black/delIP",
                type: "post",
                data: JSON.stringify(del_data),
                contentType: 'application/json;charset=utf-8',
                dataType: "json"
            }).done(function (msg) {
                if (msg.code === 200) {
                    render_black_ip_table(netcard);
                }else{
                    layer.msg(msg.msg)
                }
            }).fail(function (e) {
                layer.msg("error")
            })
        })

        // 添加IP黑名单
        $("#add-black-ip").click(function(){
            //例子2
            layer.prompt({
                formType: 3,
                value: '',
                title: '请输入IP地址',
                // area: ['800px', '350px'] //自定义文本域宽高
            }, function(value, index, elem){
                // alert(value); //得到value
                if (!validateIP(value)){
                    layer.msg("请检查IP地址格式")
                    return
                }
                var netcard = $("#netcard-select-input").val();
                if (netcard==="default"){
                    layer.msg("请先选择一张网卡")
                    return
                }
                $.ajax({
                    url:"/xdp/black/setIP",
                    type:"post",
                    data:JSON.stringify({
                        iface: netcard,
                        blackIpList:[
                            value,
                        ]
                    }),
                    dataType:"json",
                    contentType: "application/json;charset=utf-8",
                    async: false,
                }).done(function(msg){
                    if (msg.code!==200){
                        layer.msg(msg.msg)
                    }
                }).fail(function(e){
                    layer.msg("error")
                })
                layer.close(index);
                render_black_ip_table(netcard);
            });
        })

    })

})

