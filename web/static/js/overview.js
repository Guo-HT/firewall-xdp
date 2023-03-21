$(function () {
    layui.use(['layer'], function () {
        var layer = layui.layer;

        var is_first_requests = true
        var first_timestamp
        var first_server_runtime
        var line_chart_netcard_index = 0
        var rx, tx

        var clock_chart = echarts.init(document.getElementById("clock"));
        var clock_static_chart = echarts.init(document.getElementById("clock_static"));
        var cpu_chart = echarts.init(document.getElementById("cpu"));
        var mem_chart = echarts.init(document.getElementById("mem"));
        var disk_chart = echarts.init(document.getElementById("disk"));
        var temp_chart = echarts.init(document.getElementById("temp"));
        var netcard_speed_chart = echarts.init(document.getElementById("netcard-speed-charts"));

        var _speed_init_rx = []
        var _speed_init_tx = []
        for (var i = 0; i <= 60; i++) {
            // _speed_init.push(0)
            // _speed_init_rx.push(Math.random() * 100)
            // _speed_init_tx.push(Math.random() * 100)
            _speed_init_rx.push(0)
            _speed_init_tx.push(0)
        }
        var xAxis_data_init = []
        for (var j = 0; j <= 60; j++) {
            xAxis_data_init.push(j)
        }

        var clock_option = {
            series: [
                {
                    name: '小时',
                    type: 'gauge',
                    min: 0,
                    max: 12,
                    splitNumber: 12,
                    radius: '60%',
                    startAngle: 90,
                    endAngle: -269.999,
                    animation: 0,
                    axisLine: {
                        // 坐标轴线
                        lineStyle: {
                            // 属性lineStyle控制线条样式
                            color: [[1, '#C0911F']], width: 0
                        }
                    },
                    axisLabel: {
                        show: 0, fontSize: 0
                    },
                    axisTick: {
                        // 坐标轴小标记
                        show: false
                    },
                    splitLine: {
                        // 分隔线
                        show: false
                    },
                    pointer: {
                        // 分隔线
                        shadowColor: '#fff', //默认透明
                        shadowBlur: 10, length: '50%', width: 3, color: '#ffffff', itemStyle: {
                            color: '#FFFFFF'
                        }
                    },
                    detail: {
                        show: false
                    },
                    data: [{value: 10, name: ''}]
                },
                {
                    name: '分钟',
                    type: 'gauge',
                    min: 0,
                    max: 60,
                    splitNumber: 12,
                    radius: '60%',
                    startAngle: 90,
                    endAngle: -269.999,
                    animation: 0,
                    axisLine: {
                        // 坐标轴线
                        lineStyle: {
                            // 属性lineStyle控制线条样式
                            color: [[1, '#C0911F']], width: 0
                        }
                    },
                    axisLabel: {
                        show: false
                    },
                    axisTick: {
                        // 坐标轴小标记
                        show: false
                    },
                    splitLine: {
                        // 分隔线
                        show: 0, fontSize: 0
                    },
                    pointer: {
                        // 分隔线
                        shadowColor: '#fff', //默认透明
                        shadowBlur: 10, length: '60%', width: 2
                    },
                    detail: {
                        show: false
                    },
                    data: [{value: 10, name: ''}]
                },
                {
                    name: '秒',
                    type: 'gauge',
                    min: 0,
                    max: 60,
                    splitNumber: 12,
                    radius: '60%',
                    startAngle: 90,
                    endAngle: -269.999,
                    animation: 0,
                    axisLine: {
                        // 坐标轴线
                        lineStyle: {
                            // 属性lineStyle控制线条样式
                            color: [[1, '#000000']], width: 2, shadowColor: '#fff', //默认透明
                            shadowBlur: 10
                        }
                    },
                    axisLabel: {
                        // 坐标轴小标记
                        // formatter: function(v) {
                        //   switch (v + '') {
                        //     case '0':
                        //       return ''
                        //     default:
                        //       return v
                        //   }
                        // },
                        show: 0, textStyle: {
                            // 属性lineStyle控制线条样式
                            fontWeight: 'bolder', color: '#fff', shadowColor: '#fff', //默认透明
                            shadowBlur: 10
                        }
                    },
                    axisTick: {
                        // 坐标轴小标记
                        length: 1, // 属性length控制线长
                        lineStyle: {
                            // 属性lineStyle控制线条样式
                            color: 'auto', shadowColor: 'rgba(0, 0, 0, 0.3)', //默认透明
                            shadowBlur: 5
                        }
                    },
                    splitLine: {
                        // 分隔线
                        // show:false,
                        length: 3, // 属性length控制线长
                        lineStyle: {
                            // 属性lineStyle（详见lineStyle）控制线条样式
                            width: 3, color: '#000000', shadowColor: 'rgba(0, 0, 0, 0.3)', //默认透明
                            shadowBlur: 3,
                        }
                    },
                    pointer: {
                        // 分隔线
                        shadowColor: '#fff', //默认透明
                        shadowBlur: 10, length: '80%', width: 1
                    },
                    title: {
                        textStyle: {
                            // 其余属性默认使用全局文本样式，详见TEXTSTYLE
                            fontWeight: 'bolder', fontSize: 20, fontStyle: 'italic', color: '#fff', shadowColor: '#fff', //默认透明
                            shadowBlur: 10
                        }
                    },
                    detail: {
                        show: false
                    },
                    data: [{value: 30, name: ''}]
                }]
        }

        var template_gauge_option = {
            series: [
                {
                    type: 'gauge',
                    startAngle: 200,
                    endAngle: -20,
                    min: 0,
                    max: 100,
                    splitNumber: 10,
                    center: ["50%", "55%"], // 仪表位置
                    axisLine: {
                        lineStyle: {
                            width: 6,
                            color: [
                                [0.25, '#FF6E76'],
                                [0.5, '#FDDD60'],
                                [0.75, '#58D9F9'],
                                [1, '#7CFFB2']
                            ]
                        }
                    },
                    pointer: {
                        icon: 'path://M12.8,0.7l12,40.1H0.7L12.8,0.7z',
                        length: '62%',
                        width: 4,
                        offsetCenter: [0, '-60%'],
                        itemStyle: {
                            color: 'auto'
                        }
                    },
                    axisTick: {
                        length: 8,
                        lineStyle: {
                            color: 'auto',
                            width: 2
                        }
                    },
                    splitLine: {
                        length: 14,
                        lineStyle: {
                            color: 'auto',
                            width: 3
                        }
                    },
                    axisLabel: {
                        show: false,
                        color: '#464646',
                        fontSize: 20,
                        distance: -60,
                        formatter: function (value) {
                            return '';
                        }
                    },
                    title: {
                        offsetCenter: [0, '48%'],
                        fontSize: 14
                    },
                    detail: {
                        fontSize: 20,
                        offsetCenter: [0, '-30%'],
                        valueAnimation: true,
                        formatter: function (value) {
                            return Math.round(value) + '%';
                        },
                        color: '#000000',
                    },
                    data: [
                        {
                            value: 0,
                            name: ''
                        }
                    ]
                }
            ]
        };

        var net_speed_option = {
            title: {
                // text: 'Stacked Line'
            },
            tooltip: {
                trigger: 'axis',
            },
            legend: {
                data: ['接收速率', '发送速率'],
                x: 'right',      //可设定图例在左、右、居中
                y: 'top',     //可设定图例在上、下、居中
            },
            grid: {
                left: '1%',
                right: '1%',
                bottom: '5%',
                top: '10%',
                containLabel: true
            },

            xAxis: {
                name: 's',
                type: 'category',
                boundaryGap: false,
                data: xAxis_data_init,
                interval: 5, // 步长
                min: 0, // 起始
                max: 60 // 终止
            },
            axisLabel: {
                color: "#7C7D86"
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: '接收速率',
                    type: 'line',
                    // stack: '接收速率',
                    // data: [120, 132, 101, 134, 90, 230, 210]
                    data: _speed_init_rx
                },
                {
                    name: '发送速率',
                    type: 'line',
                    // stack: '发送速率',
                    // data: [220, 182, 191, 234, 290, 330, 310]
                    data: _speed_init_tx
                }
            ]
        };

        var cpu_gauge_option = JSON.parse(JSON.stringify(template_gauge_option))
        cpu_gauge_option.series[0].axisLine.lineStyle.color = [
            [0.25, '#FF6E76'],
            [0.5, '#FDDD60'],
            [0.75, '#58D9F9'],
            [1, '#7CFFB2']
        ]
        cpu_gauge_option.series[0].detail.formatter = function (value) {
            return Math.round(value) + '%';
        }
        cpu_gauge_option.series[0].data[0].name = "CPU占用率 ( % )"

        var mem_gauge_option = JSON.parse(JSON.stringify(template_gauge_option))
        mem_gauge_option.series[0].axisLine.lineStyle.color = [
            [0.25, '#FF6E76'],
            [0.5, '#FDDD60'],
            [0.75, '#58D9F9'],
            [1, '#7CFFB2']
        ]
        mem_gauge_option.series[0].detail.formatter = function (value) {
            return Math.round(value) + '%';
        }
        mem_gauge_option.series[0].data[0].name = "内存使用率 ( % )"

        var disk_gauge_option = JSON.parse(JSON.stringify(template_gauge_option))
        disk_gauge_option.series[0].axisLine.lineStyle.color = [
            [0.25, '#FF6E76'],
            [0.5, '#FDDD60'],
            [0.75, '#58D9F9'],
            [1, '#7CFFB2']
        ]
        disk_gauge_option.series[0].detail.formatter = function (value) {
            return Math.round(value) + '%';
        }
        disk_gauge_option.series[0].data[0].name = "磁盘使用率 ( % )"

        var temp_gauge_option = JSON.parse(JSON.stringify(template_gauge_option))
        temp_gauge_option.series[0].axisLine.lineStyle.color = [
            [0.25, '#FF6E76'],
            [0.5, '#FDDD60'],
            [0.75, '#58D9F9'],
            [1, '#7CFFB2']
        ]
        temp_gauge_option.series[0].detail.formatter = function (value) {
            return Math.round(value) + '℃';
        }
        temp_gauge_option.series[0].data[0].name = "温度 ( ℃ )"

        // 表
        clock_option && clock_chart.setOption(clock_option);
        clock_option && clock_static_chart.setOption(clock_option);
        // 仪表盘
        cpu_gauge_option && cpu_chart.setOption(cpu_gauge_option);
        mem_gauge_option && mem_chart.setOption(mem_gauge_option);
        disk_gauge_option && disk_chart.setOption(disk_gauge_option);
        temp_gauge_option && temp_chart.setOption(temp_gauge_option);
        // 折线图
        net_speed_option && netcard_speed_chart.setOption(net_speed_option);

        get_system_status()
        setInterval(get_system_status, 1000)

        function get_system_status() {
            $.ajax({
                url: "/status/overview",
                type: "get",
                dataType: "json",
            }).done(function (msg) {
                // console.log(msg)
                if (msg.code === 200) {
                    var data = msg.data
                    var Y, M, D, h, m, s, run_h, run_m, run_s
                    var time_array, runtime_array
                    if (is_first_requests === true) {
                        // 第一次请求页面
                        first_timestamp = data.server_time
                        first_server_runtime = data.system_runtime
                        // console.log("第一次")
                        is_first_requests = false
                    } else {
                        // 后续定时器刷新
                        first_timestamp = first_timestamp + 1
                        first_server_runtime = first_server_runtime+1
                        // console.log("第n次")
                    }
                    time_array = parse_time(first_timestamp)
                    Y = time_array[0]
                    M = time_array[1] + 1
                    D = time_array[2]
                    h = time_array[3]
                    m = time_array[4]
                    s = time_array[5]
                    // console.log(Y, M, D, h, m, s)
                    clock_option.series[0].data[0].value = (h % 12 + m / 60).toFixed(1)
                    clock_option.series[1].data[0].value = m.toFixed(1)
                    clock_option.series[2].data[0].value = s.toFixed(1)

                    runtime_array = parse_runtime_string(first_server_runtime)
                    run_h = runtime_array[0]
                    run_m = runtime_array[1]
                    run_s = runtime_array[2]
                    $("#running-time").text(run_h+"小时"+run_m+"分"+run_s+"秒")

                    cpu_gauge_option.series[0].data[0].value = data.cpu_percent;
                    mem_gauge_option.series[0].data[0].value = data.mem_percent;
                    disk_gauge_option.series[0].data[0].value = data.disk_percent;
                    temp_gauge_option.series[0].data[0].value = data.temperature;

                    $("#system-time").text(time_array[6])
                    clock_option && clock_chart.setOption(clock_option);
                    cpu_gauge_option && cpu_chart.setOption(cpu_gauge_option);
                    mem_gauge_option && mem_chart.setOption(mem_gauge_option);
                    disk_gauge_option && disk_chart.setOption(disk_gauge_option);
                    temp_gauge_option && temp_chart.setOption(temp_gauge_option);

                    var html = network_status_html(data.speed_io)
                    $("#net-io-list").html(html)

                    var netio_this = $(".netio-this");
                    if (netio_this.length!==0){
                        // 没有点击过，默认第一个
                        line_chart_netcard_index = 0
                    }
                    // 点击过，获取索引值
                    tx = data.speed_io[line_chart_netcard_index].send_bytes_speed/1000;
                    rx = data.speed_io[line_chart_netcard_index].recv_bytes_speed/1000;
                    _speed_init_rx.shift()
                    _speed_init_rx.push(rx)
                    _speed_init_tx.shift()
                    _speed_init_tx.push(tx)
                    net_speed_option && netcard_speed_chart.setOption(net_speed_option);

                } else {
                    layer.msg(msg.msg)
                }
            }).fail(function (e) {
                layer.msg("error")
            })
        }

        // 解析时间戳
        function parse_time(timestamp) {
            var server_time = moment(moment(timestamp * 1000).format());
            return [server_time.year(), server_time.month(), server_time.date(), server_time.hours(), server_time.minutes(), server_time.seconds(), moment(timestamp * 1000).format("YYYY-MM-DD hh:mm:ss")]
        }

        // 秒 转化成 时\分\秒
        function parse_runtime_string(runtime) {
            const h = parseInt(runtime / 3600)
            const m = parseInt(runtime / 60 % 60)
            const s = Math.ceil(runtime % 60)
            return [h, m, s]
        }

        function network_status_html(data){
            var html = "";
            for(var i=0;i<data.length;i++){
                html += ("<li class=\"netcard-li layui-clear\"><div class=\"layui-col-md5\"><div class=\"netcard\"><div class=\"netcard-icon-svg\"></div></div>\n" +
                    "<div class=\"netcard netcard-name\">"+data[i].name+"</div></div><div class=\"layui-col-md7\">\n" +
                    "<div><span class=\"netcard-rxtx\">发送速率:</span>&nbsp;<span title='KB/s'>"+data[i].send_bytes_speed/1000+"</span></div>\n" +
                    "<div><span class=\"netcard-rxtx\">接受速率:</span>&nbsp;<span title='KB/s'>"+data[i].recv_bytes_speed/1000+"</span></div></div></li>")
            }
            return html
        }
        // netio-this
        $("#net-io-list").on("click", ".netcard-li", function(){
            line_chart_netcard_index = $(this).index()
            for (var i = 0; i <= 60; i++) {
                // _speed_init.push(0)
                // _speed_init_rx.push(Math.random() * 100)
                // _speed_init_tx.push(Math.random() * 100)
                _speed_init_rx.shift()
                _speed_init_rx.push(0)
                _speed_init_tx.shift()
                _speed_init_tx.push(0)
            }
            // console.log($(this).index())
        })



    })
})