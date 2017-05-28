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
            return res.rows;
        },
        onClickRow:function (row, $element, field) {
            $('#keys_body_row').show();
            $('#title_keys_make').text("编辑");
            $('#oldkey').val(row.key);
            $('#oldttl').val(row.ttl);
            $('#oldtype').val(row.type);
            $('#key_type_select').val(row.type);
            init_key_form();
            $('#key_type_select').attr('disabled','disabled');
            $('#key_form_body').show();
            $('#key_table_body').show();
            $('#key_name').val(row.key);
            $('#key_name').attr('readonly','readonly');
            $('#key_value_table').bootstrapTable('refreshOptions',keyvaluetableinit())
        }
    }
}
$('#key_value_table').bootstrapTable(keyvaluetableinit());
function keyvaluetableinit() {
    var tabledata={
        striped:true,
        uniqueId:"id",
        // cardView:true,
        pagination:true,
        pageNumber:1,
        pageSize:10,
        pageList:[10,100,1000,10000,100000],
        sidePagination:'client',
        // clickToSelect:true,
        url:'/keydata',
        queryParams:function(params) {
            params["key"]=$('#oldkey').val();
            params["type"]=$('#oldtype').val();
            params["redis"]=$('#redis_select').val();
            params["redis_db"]=$('#redis_db_select').val();
            return params;
        },
        method:'post',
        responseHandler:function(res) {
            return res.rows;
        },
    };
    switch ($('#oldtype').val())
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
                field: 'index',
                title: '索引'
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
    return tabledata
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
    redis_db_select();
    $('#keystable').bootstrapTable("refresh",{});
});

$('#redis_db_select').change(function () {
    $('#keystable').bootstrapTable("refresh",{});
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
    $('#keys_body_row').show();
    $('#title_keys_make').text("新建KEY");
    $('#key_form_body').show();
    $('#key_table_body').hide();
    init_key_form();

});

function init_key_form() {
    var key_type_select=$('#key_type_select').val();
    switch (key_type_select)
    {
        case "string":
            init_key_str_form();
            break;
        case "list":
            init_key_str_form();
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
    empty_form();
    $('#key_name_group').show();
    $('#key_val_group').show();
}

function init_key_hash_form() {
    empty_form();
    $('#key_name_group').show();
    $('#key_field_name_group').show();
    $('#key_val_group').show();
}

function init_key_zset_form() {
    empty_form();
    $('#key_name_group').show();
    $('#key_score_group').show();
    $('#key_val_group').show();
}

function empty_form() {
    $('#key_save_wran').hide();
    $('#key_type_select').removeAttr('disabled');
    $('#key_name').val("");
    $('#key_name').removeAttr("readonly");
    $('#key_field_name').val("");
    $('#key_val').val("");
    $('#key_field_name_group').hide();
    $('#key_name_group').hide();
    $('#key_val_group').hide();
    $('#key_score_group').hide();
}

$('#key_type_select').change(function () {
    init_key_form()
});


$('#key_cancel').click(function () {
    if ($('#title_keys_make').text()=="新建KEY"){
        return
    }
    $('#key_form_body').hide();
    $('#key_table_body').show();
});


$('#key_save').click(function () {
    var type=$('#key_type_select').val(),
        key=$.trim($('#key_name').val()),
        val=$.trim($('#key_val').val()),
        score=$.trim($('#key_score').val()),
        field=$.trim($('#key_field_name').val());
    if (key==""&&val==""){
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
        data:JSON.stringify({"type":type,"key":key,"val":val,"field":field,"score":score,"redis":$('#redis_select').val(),"redis_db":$('#redis_db_select').val()}),
        // traditional:true,
        contentType: "application/json",
        dataType:'json',
        success:function (res) {
            $('#keystable').bootstrapTable("refresh",{});
        }
    });
});
