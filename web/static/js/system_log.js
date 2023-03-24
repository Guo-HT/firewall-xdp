$(function () {
    layui.use(['element', 'table', "form", "layer", "upload", "laydate"], function () {
        const element = layui.element;
        const table = layui.table;
        const form = layui.form;
        const layer = layui.layer;
        const upload = layui.upload;
        const laydate = layui.laydate;

        laydate.render({
            elem: "#sys-log-end-time",
            type: "datetime"
        })
        laydate.render({
            elem: "#sys-log-start-time",
            type: "datetime"
        })

        var search_content = ""

        render_sys_log_table(null,"desc", 0, 9999999999)

        function render_sys_log_table(init_sort, sort, start_time, end_time) {
            table.render({
                elem: '#sys-log-table'
                , url: '/log/search' //数据接口
                , page: { //开启分页
                    limit: 10,
                    limits: [10, 20, 30, 50, 100],
                },
                initSort: {
                    field: "time",
                    type: init_sort
                },
                where: {
                    "search": search_content,
                    "start_time": start_time,
                    "end_time": end_time,
                    "sort": sort,
                },
                request: {
                    pageName: 'page_no' //页码的参数名称，默认：page
                    , limitName: 'page_size' //每页数据量的参数名，默认：limit
                }, parseData: function (res) {
                    var interface_data = []
                    if (res.code === 200) {
                        if (res.data.log === null || res.data.log.length === 0) {
                            return {
                                "code": res.code === 200 ? 0 : -1,
                                "msg": res.msg,
                                "count": res.data.total,
                                "data": [{
                                    "ip": "-",
                                    "username": "-",
                                    "option": "-",
                                    "option_result": "-",
                                    "create_at": "-",
                                }]
                            };
                        }
                        for (var i = 0; i < res.data.log.length; i++) {
                            interface_data.push({
                                "ip": res.data.log[i].ip,
                                "username": res.data.log[i].username,
                                "option": res.data.log[i].option,
                                "option_result": res.data.log[i].opt_result === true ? "成功" : "失败",
                                "time": moment(res.data.log[i].time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                            })
                        }
                        return {
                            "code": res.code === 200 ? 0 : -1,
                            "msg": res.msg,
                            "count": res.data.total,
                            "data": interface_data,
                        }
                    }
                }
                , cols: [[ //表头
                    // {field: 'id', title: 'ID', width:80, sort: true, fixed: 'left'}
                    {field: 'ip', title: '源IP', width: "15%", align:"center"}
                    , {field: 'username', title: '用户', width: "15%", align:"center"}
                    , {field: 'option', title: '操作详情', width: "40%", align:"center"}
                    , {field: 'option_result', title: '结果', width: "10%", align:"center"}
                    , {field: 'time', title: '时间', width: "20%", align:"center", sort: true}
                ]]
            });
        }

        table.on("sort(sys-log-table)", function (obj) {
            var sort_type = obj.type === null ? "desc" : obj.type
            var start_timestamp = moment($("#sys-log-start-time").val()).unix()
            var end_timestamp = moment($("#sys-log-end-time").val()).unix()
            render_sys_log_table(obj.type, sort_type, start_timestamp, end_timestamp)
        })

        $("#log-search-btn").click(function () {
            search_content = $("#sys-log-search").val()
            var start_timestamp = moment($("#sys-log-start-time").val()).unix()
            var end_timestamp = moment($("#sys-log-end-time").val()).unix()
            render_sys_log_table("desc", "desc", start_timestamp, end_timestamp)
        })
    })
})