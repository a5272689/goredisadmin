/**
 * Created by root on 5/18/17.
 */
redis_db_select();
$('#keystable').bootstrapTable(tableinit());

function tableinit() {
    return {
        striped:true,
        uniqueId:"rowid",
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
            params["keys"]=$.trim($('#keys_form').val());
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
            var newdata=[];
            for (var rowid in res.rows){
                var newrow=res.rows[rowid];
                newrow["rowid"]=Number(rowid);
                newdata.push(newrow)
            }
            return newdata;
        },
        onClickRow:function (row, $element, field) {
            $('#key_row_id').val(row.rowid);
            writekeyinfo()
        }
    }
}

function writekeyinfo() {
    var row=$('#keystable').bootstrapTable("getRowByUniqueId",$('#key_row_id').val());
    var senddata={"key":row.key,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()}
    $.ajax({
        url:"/keydata",
        type: "post",
        data:JSON.stringify(senddata),
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            empty_form();
            $('#title_keys_make').text("编辑");
            $('#key_type_select').val(res.type);
            $('#key_type_select').attr('disabled','disabled');
            init_key_form();
            $('#keys_body_row').show();
            $("#key_info_table").show();
            $('#key_table_body').show();
            $('#oldtonew_key_name').val(row.key);
            $('#oldkey').text(row.key);
            $('#oldttl').text(res.ttl);
            $('#oldtype').text(res.type);
            $('#key_name').val(row.key);
            $('#key_name').attr('readonly','readonly');
            $('#rename_key_name_group').show();
            $('#key_value_table').bootstrapTable('refreshOptions',keyvaluetableinit(res));
            console.log(res)
        }
    });
}



$('#key_value_table').bootstrapTable({});
function keyvaluetableinit(res) {
    var newrow=[];
    for (var index in res.rows){
        var tmpdata=res.rows[index];
        tmpdata['rowid']=Number(index);
        newrow.push(tmpdata)
    }
    var tabledata={
        striped:true,
        uniqueId:"rowid",
        // cardView:true,
        pagination:true,
        pageNumber:1,
        pageSize:10,
        pageList:[10,100,1000,10000,100000],
        sidePagination:'client',
        data:res.rows,
    };
    switch (res.type)
    {
        case "string":
            tabledata.columns=[{
                field: 'val',
                title: '值'
            }];
            break;
        case "list":
            tabledata.columns=[{
                field: 'index',
                title: '索引'
            },{
                field: 'val',
                title: '值'
            }];
            break;
        case "set":
            tabledata.columns=[{
                field: 'index',
                title: '索引'
            },{
                field: 'val',
                title: '值'
            }];
            break;
        case "zset":
            tabledata.columns=[{
                field: 'score',
                title: '权重值'
            },{
                field: 'val',
                title: '值'
            }];
            break;
        case "hash":
            tabledata.columns=[{
                field: 'field',
                title: 'field'
            },{
                field: 'val',
                title: '值'
            }];
            break;
    }
    if (res.type=="set"){
        tabledata.columns.push({
            title: '操作',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change='<button type="button" class="btn btn-danger btn-xs" onclick="delkey_val('+index+')">删除</button>';
                return change
            }});
    }else {
        tabledata.columns.push({
            title: '操作',
            align:'center',
            valign: 'middle',
            formatter:function (value,row,index) {
                var change='<button type="button" class="btn btn-primary btn-xs" onclick="editkey_val('+index+')">编辑</button>';
                change+='<button type="button" class="btn btn-danger btn-xs" onclick="delkey_val('+index+')">删除</button>';
                return change
            }});
    }

    return tabledata
}

// 监控keys输入框的输入，如果是回车，则刷新下方的key表
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
    redis_db_select();
    $('#keystable').bootstrapTable("refresh",{});
    $('#keys_body_row').hide();
});

$('#redis_db_select').change(function () {
    $('#keystable').bootstrapTable("refresh",{});
    $('#keys_body_row').hide();
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

$('#newkey').click(function () {
    empty_form();
    $('#title_keys_make').text("新建KEY");
    $('#keys_body_row').show();
    init_key_form();
});

function init_key_form() {
    var key_type=$('#key_type_select').val();
    $('#key_field_name').val("");
    $('#key_score').val("");
    switch (key_type)
    {
        case "string":
            init_key_str_form();
            break;
        case "list":
            init_key_list_form();
            break;
        case "set":
            init_key_str_form();
            break;
        case "zset":
            init_key_zset_form();
            break;
        case "hash":
            init_key_hash_form();
            break;
    }
}


function init_key_str_form() {
    $('#key_score_group').hide();
    $('#key_field_name_group').hide();
    $('#key_index_group').hide();
}

function init_key_hash_form() {
    $('#key_field_name_group').show();
    $('#key_score_group').hide();
    $('#key_index_group').hide();
}

function init_key_zset_form() {
    $('#key_score_group').show();
    $('#key_field_name_group').hide();
    $('#key_index_group').hide();
}
function init_key_list_form() {
    $('#key_score_group').hide();
    $('#key_field_name_group').hide();
    if ($('#title_keys_make').text()=="新建KEY"){
        $('#key_field_name_group').hide();
    }else {
        $('#key_index_group').show();
    }

}

function empty_form() {

    $('#key_type_select').removeAttr('disabled');
    $('#key_name').val("");
    $('#key_name').removeAttr("readonly");
    $('#key_row_id').val("");
    $('#oldkey').text("");
    $('#oldttl').text("");
    $('#oldtype').text("");
    $('#key_field_name').val("");
    $('#key_val').val("");
    $('#key_score').val("");
    $('#key_index').val("");
    $('#key_save_wran').hide();
    $("#key_info_table").hide();
    $('#rename_key_name_group').hide();
    $('#key_field_name_group').hide();
    $('#key_index_group').hide();
    $('#key_score_group').hide();
    $('#key_table_body').hide();

}

$('#key_type_select').change(function () {
    init_key_form()
});


$('#key_cancel').click(function () {
    $('#keys_body_row').hide();
});


$('#key_save').click(function () {
    var type=$('#key_type_select').val(),
        key=$.trim($('#key_name').val()),
        val=$.trim($('#key_val').val()),
        score=$.trim($('#key_score').val()),
        index=$.trim($('#key_index').val()),
        field=$.trim($('#key_field_name').val());
    if (key==""||val==""){
        $('#key_save_wran').show()
    }
    if (type=="hash"){
        if (field==""){
            $('#key_save_wran').show()
        }
    }
    if (type=="zset"){
        if (score==""){
            $('#key_save_wran').show()
        }
    }
    $.ajax({
        url:"/keysave",
        type: "post",
        data:JSON.stringify({"type":type,"key":key,"val":val,"field":field,"index":index,"score":score,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()}),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
            $('#keys_body_row').hide();
        }
    });
});

$('#rename_key_name').click(function () {
    var key=$.trim($('#key_name').val()),
        newkey=$.trim($('#oldtonew_key_name').val());
    if (newkey==""){
        $('#key_save_wran').show()
    }
    $.ajax({
        url:"/keyrename",
        type: "post",
        data:JSON.stringify({"key":key,"newkey":newkey,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()}),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
            $('#keys_body_row').hide();
        }
    });
});


function delkey_val(index) {
    var row=$('#key_value_table').bootstrapTable("getRowByUniqueId",index);
    row["key"]=$.trim($('#key_name').val());
    row["redis"]=$('#redis_select').val();
    row["redis_db"]=$('#redis_db_select').val();
    row["type"]=$('#key_type_select').val();
    $.ajax({
        url:"/keyvaldel",
        type: "post",
        data:JSON.stringify(row),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
            $('#keys_body_row').hide();
        }
    });
}

function editkey_val(index) {
    var row=$('#key_value_table').bootstrapTable("getRowByUniqueId",index);
    console.log(row);
    $('#key_val').val(row["val"]);
    switch ($('#key_type_select').val())
    {
        case "list":
            $('#key_index').val(row["index"]);
            break;
        case "zset":
            $('#key_score').val(row["score"]);
            break;
        case "hash":
            $('#key_field_name').val(row["field"]);
            break;
    }
}