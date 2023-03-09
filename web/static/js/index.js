window.onload = function(){
    layui.use(['element','jquery'], function(){
        var element = layui.element;
        var $ = layui.jquery;


        $(".my-page-nav").click(function(){
            var target = $(this).attr("goto")
            goto_page(target)
        })

        // 刷新页面时，获取最后一级，直接跳转
        goto_page(get_last_url())

        // 跳转
        function goto_page(target){
            history.pushState({}, null, location.origin + location.pathname +"#/"+target);
            $("#my-main-iframe").attr("src", "/web/static/page/"+target+".html")
        }

        // 获取最后一级uri
        function get_last_url(){
            var url = window.location.href;
            var index = url.lastIndexOf("/");
            var str = url.substring(index + 1)
            return str
        }

    });
}




