/**
 * Created by root on 5/12/17.
 */

$('#loginbutton').click(function () {
    var username=$("#login-username").val();
    var passwd=$("#login-passwd").val();
    var $loginwarn=$('#loginwarn')
    $.ajax({
        url:"/loginauth",
        type: "post",
        data:{"username":username,"passwd":passwd},
        traditional:true,
        dataType:"json",
        success:function (res) {
            if (res.result){
                location.href="/"
            }else {

                $loginwarn.text(res.info);
                $loginwarn.css("display","")
            }
        },
        error:function () {
            $loginwarn.text("请求验证失败！！！");
            $loginwarn.css("display","")
        }
    })
});
