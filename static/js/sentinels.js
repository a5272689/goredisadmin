/**
 * Created by root on 5/18/17.
 */

$('#sentinelstable').bootstrapTable(tableinit());

function tableinit() {
    return {
        // height:table_height(),
        striped:true,
        search:true,
        // data:tabledata,
        uniqueId:"id",
        searchAlign:'right',
        pagination:true,
        pageNumber:1,
        pageSize:10,
        pageList:[10,15,20,30,50],
        sidePagination:'client',
        // clickToSelect:true,
        url:'/sentinelsdata',

        method:'post',
        // queryParams:function(params) {
        //     var data ={};
        //     // var data =  {
        //     //     rows:params.limit,
        //     //     page:Math.ceil(params.offset/params.limit+1) || 1,
        //     // };
        //     // if (params.sort){
        //     //     data['sort']=params.sort;
        //     //     data['order']=params.order;
        //     // };
        //     // var seach_name=$('#seach_name').val();
        //     // if (seach_name) {data['seach_name']=seach_name;};
        //     console.log(params);
        //     return data;
        // },
        columns: [{
            checkbox:true
        }, {
            field: 'id',
            title: 'ID',
            align:'center',
            valign: 'middle'
        }, {
            field: 'hostname',
            title: '主机名(IP)',
            align:'center',
            valign: 'middle'
        }, {
            field: 'port',
            title: '端口',
            align:'center',
            valign: 'middle'
        }, {
            title: 'masters',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change='';
                for (var i in row['masters']){
                    // change+='<button type="button" class="btn btn-default btn-sm" onclick="delsentinel('+row['masters'][i]+')">'+row["masters"][i]+'</button>';
                    change+='<span class="label label-primary" onclick="lookredis(\''+row['masters'][i]+'\','+row['id']+')">'+row["masters"][i]+'</span>'
                }
                return change
            }
        }, {
            title: '连接状态',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change;
                if ( row["connection_status"]){
                    // change+='<button type="button" class="btn btn-default btn-sm" onclick="delsentinel('+row['masters'][i]+')">'+row["masters"][i]+'</button>';
                    change='<span class="label label-success">ON</span>'
                }else {
                    change='<span class="label label-danger">OFF</span>'
                }
                return change
            }
        }, {title:'操作',
                align:'center',
                valign: 'middle',
                formatter:function (value,row,index) {
                    var change='<button type="button" class="btn btn-primary btn-xs" onclick="writesentinel('+row['id']+')">编辑</button>';
                    change+='<button type="button" class="btn btn-danger btn-xs" onclick="delsentinel('+row['id']+')">删除</button>';
                    return change
                }
            }
        ],
        responseHandler:function(res) {
            return res.rows;
        }
        // onSearch:function (text) {
        //     var oldtabledata=$.parseJSON($.trim($("#appsjson").text()));
        //     if ($.trim(text)==""){
        //         $("#appstable").bootstrapTable('load',oldtabledata);
        //         return
        //     }
        //     var newdata=[];
        //     for (var i in oldtabledata) {
        //         for (var key in oldtabledata[i]){
        //             var tmpval=oldtabledata[i][key].toString();
        //             if (tmpval.match(text)!=null){
        //                 newdata.push(oldtabledata[i]);
        //                 break
        //             }
        //         }
        //     }
        //     $("#appstable").bootstrapTable('load',newdata);
        //     return
        // }
    }
}

$('#sentinelssavebutton').click(function () {
    var sentinelid=$.trim($('#sentinelid_form').val()),
        hostname=$.trim($('#hostname_form').val()),
        port=Number($.trim($('#port_form').val())),
        senddata={"port":port,"hostname":hostname};
    var $forminfo=$('#forminfo');
    if (sentinelid!=""){
        senddata["sentinelid"]=Number(sentinelid)
    }

    if (port==0||hostname==""){
        $forminfo.text("所有字段不能为空!!!");
        $forminfo.show();
        return
    }
    $.ajax({
        url:"/sentinelschange",
        type: "post",
        data:senddata,
        traditional:true,
        dataType:"json",
        success:function (res) {
            if (res.result){
                $('#tablerow').show();
                $('#formrow').hide();
            }else {
                $forminfo.text(res.info);
                $forminfo.show()
            }
        },
        error:function () {
            $forminfo.text("请求失败！！！");
            $forminfo.show()
        }
    });
});

$('#sentinelscancelbutton').click(function () {
    $('#tablerow').show();
    $('#formrow').hide();
});

$('#newsentinel').click(function () {
    forminit();
    $('#tablerow').hide();
    $('#formrow').show();
});

function writesentinel(id) {
    var rowdata=$('#sentinelstable').bootstrapTable('getRowByUniqueId',id);
    forminit(rowdata);
    $('#tablerow').hide();
    $('#formrow').show();
}

function forminit(data) {
    if (data!=null){
        $('#sentinelid_form').val(data.id);
        $('#hostname_form').val(data.hostname);
        $('#port_form').val(data.port);
    }else {
        $('#sentinelid_form').val("");
    }
    $('#forminfo').hide()
}

function lookredis(mastername,id) {
    console.log(mastername,id);
    console.log($('#sentinelstable').bootstrapTable('getRowByUniqueId', id));
}