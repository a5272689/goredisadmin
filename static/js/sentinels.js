/**
 * Created by root on 5/18/17.
 */

$('#sentinelstable').bootstrapTable(tableinit());

function tableinit() {
    return {
        striped:true,
        search:true,
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
            field: 'version',
            title: '版本号',
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
                    if ($('#userrole').val()=="ops") {
                        var change = '<button type="button" class="btn btn-danger btn-xs"  onclick="delsentinel(' + row['id'] + ')">删除</button>';
                        return change
                    }
                }
            }
        ],
        responseHandler:function(res) {
            console.log(res.rows);
            return res.rows;
        }
    }
}

$('#sentinelssavebutton').click(function () {
    var hostname=$.trim($('#hostname_form').val()),
        port=Number($.trim($('#port_form').val())),
        senddata={"port":port,"hostname":hostname};
    var $forminfo=$('#forminfo');
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
        // dataType:"json",
        success:function (res) {
            if (res.result){
                tablerowshow();
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

function forminit(data) {
    $('#forminfo').hide()
}

function lookredis(mastername,id) {
    var rediss=$('#sentinelstable').bootstrapTable('getRowByUniqueId', id)["master_rediss"][mastername];
    var redissstr="";
    for (var i in rediss){
        redissstr+="rediss="+rediss[i]["hostname"]+':'+rediss[i]["port"].toString()+"&"
    }
    window.location.href="/rediss?"+redissstr;
    console.log(rediss)
}


function tablerowshow() {
    $('#sentinelstable').bootstrapTable("refresh",{});
    $('#tablerow').show();
}

function delsentinel(id) {
    var info=$('#sentinelstable').bootstrapTable('getRowByUniqueId', id);
    var senddata=[{"hostname":info.hostname,"port":info.port}];
    $.ajax({
        url:"/sentinelsdel",
        type: "post",
        data:JSON.stringify(senddata),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#sentinelstable').bootstrapTable("refresh",{});
        },
    });
}

$('#delselectsentinels').click(function () {
        var infos=$('#sentinelstable').bootstrapTable('getAllSelections');
        var senddata=[];
        for (var i in infos){
            senddata.push({"hostname":infos[i].hostname,"port":infos[i].port})
        }
        $.ajax({
            url:"/sentinelsdel",
            type: "post",
            data:JSON.stringify(senddata),
            // traditional:true,
            contentType: "application/json",
            dataType:'json',
            success:function (res) {
                $('#sentinelstable').bootstrapTable("refresh",{});
            },
        });
}
);

function buttoninit() {
    if ($('#userrole').val()!="ops"){
        $('#newsentinel').attr('disabled','disabled');
        $('#delselectsentinels').attr('disabled','disabled')
    }
}
buttoninit();