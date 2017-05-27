/**
 * Created by root on 5/18/17.
 */

$('#redisstable').bootstrapTable(tableinit());

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
        url:'/redissdata',
        queryParams:function(params) {
            params["rediss"]=$("#hiddenpostdata").text();
            return params;
        },
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
            title: '角色',
            align:'center',
            valign: 'middle',
            field: 'role'
        }, {
            title: '运行天数',
            align:'center',
            valign: 'middle',
            field: 'uptime_in_days'
        }, {
            title: '占用内存含碎片(KB)',
            align:'center',
            valign: 'middle',
            field: 'used_memory_rss'
        }, {
            title: 'keys',
            align:'center',
            valign: 'middle',
            field: 'keys'
        }, {
            title: '连接',
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
        }, {
            title: '认证',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change;
                if ( row["auth_status"]){
                    // change+='<button type="button" class="btn btn-default btn-sm" onclick="delsentinel('+row['masters'][i]+')">'+row["masters"][i]+'</button>';
                    change='<span class="label label-success">YES</span>'
                }else {
                    change='<span class="label label-danger">NO</span>'
                }
                return change
            }
        }, {
            title: 'PING',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change;
                if ( row["ping_status"]){
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
                    var change='<button type="button" class="btn btn-primary btn-xs" onclick="lookredis('+row['id']+')">key操作</button>';
                    change+='<button type="button" class="btn btn-warning btn-xs" onclick="writeredis('+row['id']+')">修改密码</button>';
                    change+='<button type="button" class="btn btn-danger btn-xs" onclick="delredis('+row['id']+')">删除</button>';
                    return change
                }
            }
        ],
        responseHandler:function(res) {
            return res.rows;
        }
    }
}

$('#redisssavebutton').click(function () {
    var hostname=$.trim($('#hostname_form').val()),
        port=Number($.trim($('#port_form').val())),
        password=$.trim($('#password_form').val()),
        senddata={"port":port,"hostname":hostname,"password":password};
    var $forminfo=$('#forminfo');
    if (port==0||hostname==""){
        $forminfo.text("所有字段不能为空!!!");
        $forminfo.show();
        return
    }
    $.ajax({
        url:"/redisschange",
        type: "post",
        data:senddata,
        traditional:true,
        dataType:"json",
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

$('#redisscancelbutton').click(function () {
    tablerowshow();
    $('#formrow').hide();
});

$('#newredis').click(function () {
    forminit();
    $('#tablerow').hide();
    $('#formrow').show();
});

function writeredis(id) {
    var rowdata=$('#redisstable').bootstrapTable('getRowByUniqueId',id);
    $('#form_title').text("修改redis密码");
    $('#hostname_form_row').hide();
    $('#hostname_form').val(rowdata.hostname);
    $('#port_form_row').hide();
    $('#port_form').val(rowdata.port);
    $('#forminfo').hide();
    $('#tablerow').hide();
    $('#formrow').show();
}

function forminit(data) {
    $('#form_title').text("新建redis");
    $('#hostname_form_row').show();
    $('#port_form_row').show();
    $('#forminfo').hide()
}

function lookredis(id) {
    var redisdata=$('#redisstable').bootstrapTable('getRowByUniqueId', id);
    window.location.href="/keys?redis="+redisdata["hostname"]+':'+redisdata["port"].toString();
}


function tablerowshow() {
    $('#redisstable').bootstrapTable("refresh",{});
    $('#tablerow').show();
}

function delredis(id) {
    var info=$('#redisstable').bootstrapTable('getRowByUniqueId', id);
    var senddata=[{"hostname":info.hostname,"port":info.port}];
    $.ajax({
        url:"/redissdel",
        type: "post",
        data:JSON.stringify(senddata),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#redisstable').bootstrapTable("refresh",{});
        },
    });
}

$('#delselectrediss').click(function () {
        var infos=$('#redisstable').bootstrapTable('getAllSelections');
        var senddata=[];
        for (var i in infos){
            senddata.push({"hostname":infos[i].hostname,"port":infos[i].port})
        }
        $.ajax({
            url:"/redissdel",
            type: "post",
            data:JSON.stringify(senddata),
            // traditional:true,
            contentType: "application/json",
            dataType:'json',
            success:function (res) {
                $('#redisstable').bootstrapTable("refresh",{});
            },
        });
    }
);