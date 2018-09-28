<div class="modal fade" id="modal-addgroup" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">Add Friend</h6>
            </div>
            <div class="modal-body">
                <label>GroupName</label>
                <input id="groupname" type="text" class="form-control" placeholder="Enter GroupName...">
                
                <button onclick='doAddGroup();' type="button" class="btn btn-primary">Add</button>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<div class="modal fade" id="modal-renamegroup" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">Rename Group</h6>
            </div>
            <div class="modal-body">
                <label>OldGroupName</label>
                <input id="oldgroupname" type="text" class="form-control" disabled>

                <label>NewGroupName</label>
                <input id="newgroupname" type="text" class="form-control" placeholder="Enter NewGroupName...">
                
                <button onclick='doRenameGroup();' type="button" class="btn btn-primary">Rename</button>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<div class="modal fade" id="modal-modifycomment" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">Modify Friend Comment</h6>
            </div>
            <div class="modal-body">
                <label>Comment</label>
                <input id="modifycomment" type="text" class="form-control" placeholder="Enter Comment...">
                
                <button onclick='doModifyComment();' type="button" class="btn btn-primary">Modify</button>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<div class="modal fade" id="modal-info" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">User Info</h6>
            </div>
            <div class="modal-body">
                <label>ID</label>
                <input id="info-id" type="text" class="form-control" disabled>
                <label>Nickname</label>
                <input id="info-nickname" type="text" class="form-control">
                <label>Desc</label>
                <input id="info-desc" type="text" class="form-control">
                <label>Birthday</label>
                <input id="info-birthday" type="text" class="form-control">
                <label>Country</label>
                <input id="info-country" type="text" class="form-control">
                
                <button onclick='doUpdateAppData();' type="button" class="btn btn-primary">Save</button>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<div class="modal fade" id="modal-search" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">Search</h6>
            </div>
            <div class="modal-body">
                <div id="search-tabs">
                    <ul>
                        <li><a href="#tabs-id">Search User By Id</a></li>
                        <li><a href="#tabs-nickname">Search User By Nickname</a></li>
                        <li><a href="#tabs-roomname">Search Room</a></li>
                    </ul>
                    <div id="tabs-id">
                        <form role="form">
                            <div class="form-group">
                                <label for="search-id">ID</label>
                                <input type="text" class="form-control" id="search-id" placeholder="User ID">
                            </div>
                        </form>
                        <button onclick="doSearchUserById();" class="btn btn-primary">Search</button>
                        <div class="box box-default no-padding">
                            <li id="search-frienditem" class="hide">
                                <img src="static/dist/img/user1-128x128.jpg" alt="User Image">
                                <a class="users-list-name" href="#">Alexander Pierce</a>
                                <span class="users-list-date">Today</span>
                                <button class="btn btn-primary">Add</button>
                                <span class="join">已是好友</span>
                            </li>
                            <ul class="users-list clearfix">
                                
                            </ul>
                        </div>
                    </div>
                    <div id="tabs-nickname">
                        <form role="form">
                            <div class="form-group">
                                <label for="search-nickname">Nickname</label>
                                <input type="text" class="form-control" id="search-nickname" placeholder="User Nickname">
                            </div>
                        </form>
                        <button onclick="doSearchUserByNickname();" class="btn btn-primary">Search</button>
                        <div class="box box-default no-padding">
                            <ul class="users-list clearfix">

                            </ul>
                        </div>
                    </div>
                    <div id="tabs-roomname">
                        <form role="form">
                            <div class="form-group">
                                <label for="search-roomname">Room Name</label>
                                <input type="text" class="form-control" id="search-roomname" placeholder="Room Name">
                            </div>
                        </form>
                        <button onclick="doSearchRoom();" class="btn btn-primary">Search</button>
                        <div class="box box-default no-padding">
                            <li id="search-roomitem" class="hide">
                                <img src="static/dist/img/user1-128x128.jpg" alt="User Image">
                                <a class="users-list-name" href="#">Alexander Pierce</a>
                                <span class="users-list-date">Today</span>
                                <button class="btn btn-primary">Join</button>
                                <span class="join">已加入</span>
                            </li>
                            <ul class="users-list clearfix">
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<div class="modal fade" id="modal-createroom" style="display: none;">
    <div class="modal-dialog" style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                <span aria-hidden="true">×</span></button>
            <h6 class="modal-title">Create Room</h6>
            </div>
            <div class="modal-body">
                <label>Room Name</label>
                <input id="createroom-roomname" type="text" class="form-control">
                
                <label>Jieshao</label>
                <input id="createroom-jieshao" type="text" class="form-control">
                <label>Notice</label>
                <input id="createroom-notice" type="text" class="form-control">

                <label>Room Type</label>
                <select id="createroom-roomtype" class="form-control" onchange="if(this.value == 3){$('#createroom-password').removeClass('hide');$('#createroom-password-lb').removeClass('hide');}else {$('#createroom-password').addClass('hide');$('#createroom-password-lb').addClass('hide');}" style="width:100px">
                    <option value="1">所有人</option>
                    <option value="2">审核加入</option>
                    <option value="3">密码</option>
                </select>
                <label id="createroom-password-lb" class="hide">Password</label>
                <input id="createroom-password" type="text" class="form-control hide">
                
                <button onclick='doCreateRoom();' type="button" class="create btn btn-primary">Create</button>
                <button onclick='doSaveRoom();' type="button" class="save hide btn btn-primary">Save</button>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Close</button>
            </div>
        </div>
    <!-- /.modal-content -->
    </div>
    <!-- /.modal-dialog -->
</div>

<script>
$( function() {
    $( "#info-birthday" ).datepicker();
    $("#info-country").countrySelect();
    $( "#search-tabs" ).tabs();
} );

function clearSearchUserById() {
    $( "#tabs-id" ).find(".users-list").html("");
}
function isMyFriend(id) {
    return frienddatabyid[id] != null || (id == userdata.id);
}
function addSearchUserById(jsondata) {
    var newitem = $( "#search-frienditem" ).clone();
    newitem.attr("id", "search-frienditem-" + jsondata.id);
    newitem.removeClass("hide");

    newitem.data("user", jsondata);

    newitem.find(".users-list-name").html(jsondata.desc);
    newitem.find(".users-list-date").html(jsondata.nickname);
    
    if(isMyFriend(jsondata.id)) {
        newitem.find("button").addClass("hide");
    } else {
        newitem.find("button").click(function(){
            addFriend(jsondata.id);
        });
        newitem.find(".join").addClass("hide");
    }

    $( "#tabs-id" ).find(".users-list").append(newitem);
}

function clearSearchUserByNickname() {
    $( "#tabs-nickname" ).find(".users-list").html("");
}
function addSearchUserByNickname(jsondata) {
    var newitem = $( "#search-frienditem" ).clone();
    newitem.attr("id", "search-frienditem-" + jsondata.id);
    newitem.removeClass("hide");

    newitem.data("user", jsondata);

    newitem.find(".users-list-name").html(jsondata.desc);
    newitem.find(".users-list-date").html(jsondata.nickname);
    
    if(isMyFriend(jsondata.id)) {
        newitem.find("button").addClass("hide");
    } else {
        newitem.find("button").click(function(){
            addFriend(jsondata.id);
        });
        newitem.find(".join").addClass("hide");
    }

    $( "#tabs-nickname" ).find(".users-list").append(newitem);
}

function clearSearchRoom() {
    $( "#tabs-roomname" ).find(".users-list").html("");
}
function isRoomJoined(rid) {
    return roomdata[rid] != null;
}
function addSearchRoom(jsondata) {
    var newitem = $( "#search-roomitem" ).clone();
    newitem.attr("id", "search-roomitem-" + jsondata.rid);
    newitem.removeClass("hide");

    newitem.data("room", jsondata);

    newitem.find(".users-list-name").html(jsondata.jieshao);
    newitem.find(".users-list-date").html(jsondata.roomname);
    
    if(isRoomJoined(jsondata.rid)) {
        newitem.find("button").addClass("hide");
    } else {
        newitem.find("button").click(function(){
            reqJoinRoom(jsondata.rid);
        });
        newitem.find(".join").addClass("hide");
    }

    $( "#tabs-roomname" ).find(".users-list").append(newitem);
}

function doSearchUserById() {
    var strid = $("#search-id").val();
    if(strid == "") {
        alert("strid should not be empty!");
        return;
    }
    reqSearchUserById(strid);
}

function doSearchUserByNickname() {
    var nickname = $("#search-nickname").val();
    if(nickname == "") {
        alert("nickname should not be empty!");
        return;
    }
    reqSearchUserByNickname(nickname);
}

function doSearchRoom() {
    var roomname = $("#search-roomname").val();
    if(roomname == "") {
        alert("roomname should not be empty!");
        return;
    }
    reqSearchRoom(roomname);
}

function doAddGroup(){
    var groupname = $("#groupname").val();
    if(groupname == "") {
        alert("groupname should not be empty!");
        return;
    }
    reqCreateGroup(groupname);
    $("#modal-addgroup").modal("hide");
}

function doRenameGroup(){
    var newgroupname = $("#newgroupname").val();
    if(newgroupname == "") {
        alert("newgroupname should not be empty!");
        return;
    }
    reqRenameGroup($("#oldgroupname").val(), newgroupname);
    $("#modal-renamegroup").modal("hide");
}

function doModifyComment(){
    var modifycomment = $("#modifycomment").val();
    if(modifycomment == "") {
        alert("modifycomment should not be empty!");
        return;
    }
    modifyFriendComment($("#modifycomment").data('id'), modifycomment);
    $("#modal-modifycomment").modal("hide");
}

function doUpdateAppData(){
    var nickname = $("#info-nickname").val();
    if(nickname == "") {
        alert("nickname should not be empty!");
        return;
    }
    var desc = $("#info-desc").val();
    var birthday = $("#info-birthday").val();
    var country = $("#info-country").val();

    var udata = $("#modal-info").data("user");

    var flag = false;
    var jsondata = {};
    if(udata.nickname != nickname){
        jsondata["nickname"] = nickname;
        flag = true;
    }
    if(udata.desc != desc){
        jsondata["desc"] = desc;
        flag = true;
    }
    if(udata.birthday != birthday){
        jsondata["birthday"] = birthday;
        flag = true;
    }
    if(udata.country != country){
        jsondata["country"] = country;
        flag = true;
    }
    updateAppdata(jsondata);
    $("#modal-info").modal("hide");
}

function showAddGroupPanel() {
    $("#modal-addgroup").modal("show");
}

function showRenameGroupPanel(groupname) {
    $("#oldgroupname").val(groupname);
    $("#modal-renamegroup").modal("show");
}

function showModifyCommentPanel(idstr) {
    $("#modifycomment").data('id', idstr);
    $("#modal-modifycomment").modal("show");
}

function showUserInfoPanel(jsondata) {
    $("#modal-info .modal-title").html("User Info-<b>" + jsondata.nickname + "</b>");
    $("#info-id").val(jsondata.id);
    $("#info-nickname").val(jsondata.nickname);
    $("#info-desc").val(jsondata.desc);
    $("#info-birthday").val(jsondata.birthday);
    $("#info-country").val(jsondata.country);
    if(userdata.id != jsondata.id) {
        $("#info-nickname").attr("disabled", true);
        $("#info-desc").attr("disabled", true);
        $("#info-birthday").attr("disabled", true);
        $("#info-country").attr("disabled", true);
    } else {
        $("#info-nickname").attr("disabled", false);
        $("#info-desc").attr("disabled", false);
        $("#info-birthday").attr("disabled", false);
        $("#info-country").attr("disabled", false);
    }
    $("#info-country").countrySelect("refresh");
    $("#modal-info").data("user", jsondata);
    $("#modal-info").modal("show");
}

function showCreateRoomPanel() {
    $("#createroom-roomname").attr("disabled", false);
    $("#createroom-jieshao").attr("disabled", false);
    $("#createroom-notice").attr("disabled", false);
    $("#createroom-roomtype").attr("disabled", false);

    // $("#createroom-password").addClass("hide");
    // $('#createroom-password-lb').addClass('hide');

    $("#createroom-roomname").val("");
    $("#createroom-jieshao").val("");
    $("#createroom-notice").val("");
    $("#createroom-password").val("");
    //$("#createroom-roomtype option[index='0']").attr("selected",true);
    //$('#createroom-roomtype').get(0).selectedIndex=0;
    //$("#createroom-roomtype option:first").prop("selected", 'selected'); 
    //$("#createroom-roomtype option[value='1']").attr("selected","selected");
    $('#createroom-roomtype').val('1');
    //$("#select_id option:last")

    $("#modal-createroom").find(".save").addClass("hide");
    $("#modal-createroom").find(".create").removeClass("hide");

    $("#modal-createroom").modal("show");
}

function showRoomInfoPanel(roomdata) {
    $("#modal-createroom").data("room", roomdata);
    console.info("index="+(roomdata.roomtype - 1));

    $("#createroom-roomname").val(roomdata.roomname);
    $("#createroom-jieshao").val(roomdata.jieshao);
    $("#createroom-notice").val(roomdata.notice);
    $("#createroom-password").val("");
    //$("#createroom-roomtype option[value='"+(roomdata.roomtype)+"']").attr("selected","selected");
    $('#createroom-roomtype').val(roomdata.roomtype);
    //$("#select_id option:last")

    if(roomdata.ownerid != userdata.id){
        $("#createroom-roomname").attr("disabled", true);
        $("#createroom-jieshao").attr("disabled", true);
        $("#createroom-notice").attr("disabled", true);
        $("#createroom-roomtype").attr("disabled", true);
        $("#createroom-password").addClass("hide");
        $('#createroom-password-lb').addClass('hide');

        $("#modal-createroom").find(".save").addClass("hide");
        $("#modal-createroom").find(".create").addClass("hide");
    }else{
        if(roomdata.roomtype == 3) {
            $("#createroom-password").removeClass("hide");
            $('#createroom-password-lb').removeClass('hide');
        }
        
        $("#modal-createroom").find(".save").removeClass("hide");
        $("#modal-createroom").find(".create").addClass("hide");
    } 

    $("#modal-createroom").modal("show");
}

function doCreateRoom() {
    var roomname = $("#createroom-roomname").val();
    if(roomname == "") {
        alert("roomname should not be empty!");
        return;
    }

    var roomtype = $("#createroom-roomtype").val();
    var password = $("#createroom-password").val();
    if(roomtype == 3 && password == "") {
        alert("password should not be empty!");
        return;
    }

    var jsondata = {};
    jsondata["roomname"] = roomname;
    jsondata["roomtype"] = roomtype;
    if(roomtype == 3) {
        jsondata["password"] = password;
    }
    jsondata["jieshao"] = $("#createroom-jieshao").val();
    jsondata["notice"] = $("#createroom-notice").val();

    reqCreateRoom(jsondata);
    $("#modal-createroom").modal("hide");
}
</script>