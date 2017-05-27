/**
 * Created by root on 5/18/17.
 */
redis_db_select();
$('#keystable').bootstrapTable(tableinit());

function tableinit() {
    return {
        striped:true,
        uniqueId:"id",
        // cardView:true,
        pagination:true,
        pageNumber:1,
        pageSize:10,
        pageList:[10,100,1000,10000,100000],
        sidePagination:'client',
        // clickToSelect:true,
        url:'/keysdata',
        queryParams:function(params) {
            params["redis"]=$('#redis_select').val();
            params["keys"]=$('#keys_form').val();
            params["redis_db"]=$('#redis_db_select').val();
            return params;
        },
        method:'post',
        columns: [{
            checkbox:true
        }, {
            field: 'key',
            title: 'key'
            // align:'center',
            // valign: 'middle',
        }
        ],
        responseHandler:function(res) {
            return res.rows;
        },
        onClickRow:function (row, $element, field) {
            console.log(row)
        }
    }
}


$('#keys_form').keydown(function(event){
    if (event.keyCode==13) {
        $('#keystable').bootstrapTable("refresh",{});
    }
});

$('#keysdelbutton').click(function(event){
    var selectkeys=$('#keystable').bootstrapTable("getAllSelections",{});
    var delkeys=[];
    for (var i in selectkeys){
        delkeys.push(selectkeys[i]["key"])
    }
    var senddata={"keys":delkeys,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()};
    $.ajax({
        url:"/keysdel",
        type: "post",
        data:JSON.stringify(senddata),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
        }
    });
});

$('#keyssetttlbutton').click(function(event){
    var selectkeys=$('#keystable').bootstrapTable("getAllSelections",{});
    var delkeys=[];
    for (var i in selectkeys){
        delkeys.push(selectkeys[i]["key"])
    }
    var senddata={"keys":delkeys,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val(),"seconds":Number($.trim($('#ttlseconds').val()))};
    $.ajax({
        url:"/keysexpire",
        type: "post",
        data:JSON.stringify(senddata),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
        }
    });
});


$('#keysPersistbutton').click(function(event){
    var selectkeys=$('#keystable').bootstrapTable("getAllSelections",{});
    var delkeys=[];
    for (var i in selectkeys){
        delkeys.push(selectkeys[i]["key"])
    }
    var senddata={"keys":delkeys,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()};
    $.ajax({
        url:"/keyspersist",
        type: "post",
        data:JSON.stringify(senddata),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
        }
    });
});

$('#redis_select').change(function () {
    redis_db_select()
});

function redis_db_select() {
    var redis_name=$('#redis_select').val();
    var db_map=JSON.parse($('#db_map').text());
    var redis_db_select=$('#redis_db_select');
    redis_db_select.empty();
    for (var i in db_map[redis_name]){
        redis_db_select.append('<option value="'+db_map[redis_name][i]+'">'+db_map[redis_name][i]+'</option>')
    }
}


// $('#redisssavebutton').click(function () {
//     var hostname=$.trim($('#hostname_form').val()),
//         port=Number($.trim($('#port_form').val())),
//         password=$.trim($('#password_form').val()),
//         senddata={"port":port,"hostname":hostname,"password":password};
//     var $forminfo=$('#forminfo');
//     if (port==0||hostname==""){
//         $forminfo.text("所有字段不能为空!!!");
//         $forminfo.show();
//         return
//     }
//     $.ajax({
//         url:"/redisschange",
//         type: "post",
//         data:senddata,
//         traditional:true,
//         dataType:"json",
//         success:function (res) {
//             if (res.result){
//                 tablerowshow();
//                 $('#formrow').hide();
//             }else {
//                 $forminfo.text(res.info);
//                 $forminfo.show()
//             }
//         },
//         error:function () {
//             $forminfo.text("请求失败！！！");
//             $forminfo.show()
//         }
//     });
// });
//
// $('#redisscancelbutton').click(function () {
//     tablerowshow();
//     $('#formrow').hide();
// });
//
// $('#newredis').click(function () {
//     forminit();
//     $('#tablerow').hide();
//     $('#formrow').show();
// });
//
// function writeredis(id) {
//     var rowdata=$('#redisstable').bootstrapTable('getRowByUniqueId',id);
//     $('#form_title').text("修改redis密码");
//     $('#hostname_form_row').hide();
//     $('#hostname_form').val(rowdata.hostname);
//     $('#port_form_row').hide();
//     $('#port_form').val(rowdata.port);
//     $('#forminfo').hide();
//     $('#tablerow').hide();
//     $('#formrow').show();
// }
//
// function forminit(data) {
//     $('#form_title').text("新建redis");
//     $('#hostname_form_row').show();
//     $('#port_form_row').show();
//     $('#forminfo').hide()
// }
//
// function lookredis(id) {
//     var redisdata=$('#redisstable').bootstrapTable('getRowByUniqueId', id);
//     console.log(redisdata)
// }
//
//
// function tablerowshow() {
//     $('#redisstable').bootstrapTable("refresh",{});
//     $('#tablerow').show();
// }
//
// function delredis(id) {
//     var info=$('#redisstable').bootstrapTable('getRowByUniqueId', id);
//     var senddata=[{"hostname":info.hostname,"port":info.port}];
//     $.ajax({
//         url:"/redissdel",
//         type: "post",
//         data:JSON.stringify(senddata),
//         // traditional:true,
//         contentType: "application/json",
//         dataType:'json',
//         success:function (res) {
//             $('#redisstable').bootstrapTable("refresh",{});
//         },
//     });
// }
//
// $('#delselectrediss').click(function () {
//         var infos=$('#redisstable').bootstrapTable('getAllSelections');
//         var senddata=[];
//         for (var i in infos){
//             senddata.push({"hostname":infos[i].hostname,"port":infos[i].port})
//         }
//         $.ajax({
//             url:"/redissdel",
//             type: "post",
//             data:JSON.stringify(senddata),
//             // traditional:true,
//             contentType: "application/json",
//             dataType:'json',
//             success:function (res) {
//                 $('#redisstable').bootstrapTable("refresh",{});
//             },
//         });
//     }
// );