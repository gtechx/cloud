
console.info(Long.fromString("18446744073709551615", true, 10).toBytes());
console.info(5/2);
console.info((new Date()).getTime());
var timestamp4 = new Date((new Date()).getTime());
var ThisInt = '1529994598312'
console.info(parseInt(ThisInt))

var messagelist = {};
var userdata = {};
var frienddata = {};
var frienddatabyid = {};
var roomdata = {};
var roomuserdata = {};
var myapp = App.new();
function login(account, password, appname, zonename) {
  $.post("/webapp/chatlogin", { 'account': account, 'password': password, 'appname': appname },
    function(data) {
      console.info(data);
      var retjson = JSON.parse(data);
      if(retjson.errcode == 0) {
        myapp.onlogined = onLogined
        myapp.onloginfailed = function(errcode) {
          console.info("login failed, errcode:" + errcode)
        };
        myapp.onconnected = function(){
          //var loginjson = {};
          //loginjson.token = retjson.token;
          myapp.login(retjson.token);
        };
        myapp.onclose = function() {
          console.info("disconnect from server");
        };
        myapp.onerror = function(evt) {
          console.info("error:"+evt.data);
        };
        myapp.onkickout = function() {
          console.info("you have been kicked out from server.");
        };
        myapp.onpresence = onPresence;
        myapp.onroompresence = onRoomPresence;
        myapp.onmessage = onMessage;
        myapp.onroommessage = onRoomMessage;
        myapp.connect(retjson.serveraddr);//"127.0.0.1:9090");
      } else {
        console.error(retjson.errordesc)
      }
  });
}

function setPlatform(platform) {
  myapp.setplatform(platform);
}

function onLogined(errcode) {
  if(errcode == 0) {
    console.info("login success");
    $("#loginpanel").addClass('hide');
    myapp.requserdata("0", onMyData);
  }else{
    console.info("login failed errcode:" + errcode);
  }
}
// function onLogined(idlist) {
//   console.info(idlist);
//   if(idlist.length == 0){
//     myapp.createappdata("testnickname", onAppDataCreated);
//   } else {
//     var html = '';
//     for(var i = 0; i < idlist.length; i++) {
//       var appdataid = idlist[i];
//       console.info("appdataid:" + appdataid.toString());
//       html += '<button type="button" onclick="enterChat(\''+appdataid.toString()+'\');" style="min-width:200px;" class="btn btn-default btn-block">';
//       html += 'Sign in with ID: ' + appdataid;
//       html += '</button>';
//     }
//     $("#idlist").html(html);

//     $("#idselect").removeClass('hide');
//     $("#loginpanel").addClass('hide');
//   }
// }

function onAppDataCreated(errcode, appdataid) {
  console.info("onAppDataCreated errcode:" + errcode);
  if(errcode == 0) {
    console.info("appdataid:" + appdataid);
  }
}

function enterChat(strid) {
  console.info("strid:" + typeof(strid));
  myapp.enterchat(strid, onEnterChat);
}

function onEnterChat(errcode){
  console.info("onEnterChat errcode:" + errcode);
  if(errcode == 0) {
    myapp.requserdata("0", onMyData);
  }
}

function onMyData(errcode, jsondata) {
  console.info("onMyData errcode:" + errcode);
  if(errcode == 0) {
    console.info(jsondata);
    userdata = jsondata;
    $("#idselect").addClass('hide');
    $("#fpanel").removeClass('hide');
    $("#fpanelheader .box-title").html(jsondata.nickname);
    $("#fpanelheader .box-title").css("cursor", "hand");
    $("#fpanelheader .box-title").click(function(e){
      reqUserData(jsondata.id);
    });

    reqFriendList();
    reqPresenceList();
    myapp.reqofflinemsglist();
    reqRoomList();
  }
}

function reqUserData(idstr) {
  myapp.requserdata(idstr, onUserData);
}

function onUserData(errcode, jsondata) {
  console.info("onUserData errcode:" + errcode);
  if(errcode == 0) {
    console.info(jsondata);
    showUserInfoPanel(jsondata);
  }
}

function modifyFriendComment(idstr, comment) {
  myapp.modifyfriendcomment(idstr, comment, onModifyCommentResult)
}

function onModifyCommentResult(errcode) {
  console.info("onModifyCommentResult errcode:" + errcode);
  if(errcode == 0) {
    reqFriendList();
  }
}

function addFriend(idstr, message) {
  myapp.addfriend(idstr, message, onPresenceResult);
}

function delFriend(idstr, message) {
  myapp.delfriend(idstr, onPresenceResult);
}

function agreeFriend(idstr) {
  myapp.agreefriend(idstr, onPresenceResult);
}

function refuseFriend(idstr) {
  myapp.refusefriend(idstr, onPresenceResult);
}

function reqPresenceList() {
  myapp.reqpresencelist(onPresenceList);
}

function onPresenceList(errcode, data) {
  console.info("onPresenceList errcode:" + errcode);
  console.info("onPresenceList data length:" + data.length);
  console.info("data:" + JSON.stringify(data));
  for(var i = 0; i < data.length; i ++) {
    //console.info("onPresenceList data:" + JSON.stringify(data[i]));
    addPresence(data[i]);
    if(data[i].presencetype == PresenceType.PresenceType_UnSubscribe){
      removeFriendItemById(data[i].who);
    }
  }
}

function reqFriendList() {
  myapp.reqfriendlist(onFriendList);
}

function onFriendList(errcode, data) {
  console.info("onFriendList errcode:" + errcode);
  console.info("onFriendList data length:" + data.length);
  console.info("data:" + JSON.stringify(data));
  clearFriendList();
  frienddata = data;
  for(var groupname in data){
    //console.info("onFriendList group:" + groupname);
    createGroup(groupname);
    for(var i in data[groupname]){
      frienddatabyid[data[groupname][i].who] = data[groupname][i];
      //console.info("onFriendList item:" + JSON.stringify(data[groupname][i]));
      addFriendItem(data[groupname][i]);
    }
  }
}

function reqBlackList() {
  myapp.reqblacklist(onBlackList);
}

function onBlackList(errcode, data) {
  console.info("onBlackList errcode:" + errcode);
  console.info("onBlackList data length:" + data.length);
  // clearFriendList();
  // frienddata = data;
  // for(var group in data){
  //   console.info("onFriendList group:" + group);
  //   createGroup(group);
  //   for(var i in data[group]){
  //     console.info("onFriendList item:" + JSON.stringify(data[group][i]));
  //     addFriendItem(data[group][i]);
  //   }
  // }
}

function onPresenceResult(errcode) {
  console.info("onPresenceResult errcode:" + errcode);
  reqFriendList();
}

function onPresence(jsondata) {
  console.info("onPresence jsondata:" + JSON.stringify(jsondata));
  addPresence(jsondata);
  if(jsondata.presencetype == PresenceType.PresenceType_UnSubscribe){
    console.info("PresenceType_UnSubscribe " + jsondata.who)
    removeFriendItemById(jsondata.who);
  } else if(jsondata.presencetype == PresenceType.PresenceType_Subscribed){
    reqFriendList();
  }
}

function sendMessage(msg) {
  myapp.sendmessage(msg, onMessageResult);
}

function onMessageResult(errcode) {
  console.info("onMessageResult errcode:" + errcode);
}

function onMessage(jsondata) {
  console.info("onMessage jsondata:" + JSON.stringify(jsondata));
  var msgarray = messagelist[jsondata.who];
  if(msgarray == null){
    msgarray = new Array();
    messagelist[jsondata.who] = msgarray;
  }
  msgarray[msgarray.length] = jsondata;

  if(jsondata.from != userdata.id)
    addMessage(jsondata);
  else if(jsondata.platform != myapp.platform)
    addSendMessage(jsondata);
}

function quitChat() {
  myapp.quitchat();
  $("#loginpanel").removeClass('hide');
  $("#idselect").addClass('hide');
}

//group start
function reqCreateGroup(name) {
  myapp.creategroup(name, onGroupResult);
}

function reqDeleteGroup(name) {
  if(frienddata[name].length > 0){
    alert("can't delete group that not empty!");
    return;
  }
  myapp.deletegroup(name, onGroupResult);
}

function reqRenameGroup(oldname, newname) {
  myapp.renamegroup(oldname, newname, onGroupResult);
}

function reqMoveToGroup(idstr, name) {
  myapp.movetogroup(idstr, name, onGroupResult);
}

function onGroupResult(errcode) {
  console.info("onGroupResult errcode:" + errcode);
  reqFriendList();
}

function reqRefreshGroup(name) {
  myapp.refreshgroup(name, onRefreshGroupResult);
}

function onRefreshGroupResult(errcode, jsondata) {
  console.info("onRefreshGroupResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(jsondata));
  for(var groupname in jsondata){
    clearGroupFriendList(groupname)
    frienddata[groupname] = jsondata[groupname];
    //console.info("onRefreshGroupResult group:" + groupname);
    for(var i in jsondata[groupname]){
      //console.info("onRefreshGroupResult item:" + JSON.stringify(jsondata[groupname][i]));
      frienddatabyid[jsondata[groupname][i].who] = jsondata[groupname][i];
      addFriendItem(jsondata[groupname][i]);
    }
  }
}
//group end

//black start
function reqAddBlack(idstr) {
  myapp.addblack(idstr, onBlackResult);
}

function reqRemoveBlack(idstr) {
  myapp.removeblack(idstr, onBlackResult);
}

function onBlackResult(errcode) {
  console.info("onBlackResult errcode:" + errcode);
}
//black end

//appdata update
function updateAppdata(jsondata) {
  myapp.updateappdata(jsondata, onUpdateAppdataResult);
}

function onUpdateAppdataResult(errcode) {
  console.info("onUpdateAppdataResult errcode:" + errcode);
}
//appdata update end

//room start
function reqCreateRoom(jsondata) {
  console.info(jsondata);
  myapp.createroom(jsondata, onCreateRoomResult);
}

function onCreateRoomResult(errcode) {
  console.info("onCreateRoomResult errcode:" + errcode);
  if(errcode == 0)
  reqRoomList();
}

function reqDeleteRoom(strrid) {
  myapp.deleteroom(strrid, onDeleteRoomResult);
}

function onDeleteRoomResult(errcode) {
  console.info("onDeleteRoomResult errcode:" + errcode);
  if(errcode == 0)
  reqRoomList();
}

function reqUpdateRoomSeting(jsondata) {
  myapp.updateroomsetting(jsondata, onUpdateRoomSetingResult);
}

function onUpdateRoomSetingResult(errcode) {
  console.info("onUpdateRoomSetingResult errcode:" + errcode);
}

function reqJoinRoom(strrid, message) {
  myapp.joinroom(strrid, message, onJoinRoomResult);
}

function onJoinRoomResult(errcode) {
  console.info("onJoinRoomResult errcode:" + errcode);
  if(errcode == 0)
  reqRoomList();
}

function reqJoinRoomWithPassword(strrid, password) {
  myapp.joinroomwithpassword(strrid, password, onJoinRoomWithPasswordResult);
}

function onJoinRoomWithPasswordResult(errcode) {
  console.info("onJoinRoomWithPasswordResult errcode:" + errcode);
}

function reqQuitRoom(strrid) {
  myapp.quitroom(strrid, onQuitRoomResult);
}

function onQuitRoomResult(errcode) {
  console.info("onQuitRoomResult errcode:" + errcode);
  if(errcode == 0)
  reqRoomList();
}

function reqAgreeRoomJoin(ridstr, idstr) {
  myapp.agreeroomjoin(ridstr, idstr, onAgreeRoomJoinResult);
}

function onAgreeRoomJoinResult(errcode) {
  console.info("onAgreeRoomJoinResult errcode:" + errcode);
}

function reqRefuseRoomJoin(ridstr, idstr) {
  myapp.refuseroomjoin(ridstr, idstr, onRefuseRoomJoinResult);
}

function onRefuseRoomJoinResult(errcode) {
  console.info("onRefuseRoomJoinResult errcode:" + errcode);
}

function reqBanRoomUser(ridstr, idstr) {
  myapp.banroomuser(ridstr, idstr, onBanRoomUserResult);
}

function onBanRoomUserResult(errcode) {
  console.info("onBanRoomUserResult errcode:" + errcode);
}

function reqJinyanRoomUser(ridstr, idstr) {
  myapp.jinyanroomuser(ridstr, idstr, onJinyanRoomUserResult);
}

function onJinyanRoomUserResult(errcode) {
  console.info("onJinyanRoomUserResult errcode:" + errcode);
}

function reqUnJinyanRoomUser(ridstr, idstr) {
  myapp.unjinyanroomuser(ridstr, idstr, onUnJinyanRoomUserResult);
}

function onUnJinyanRoomUserResult(errcode) {
  console.info("onUnJinyanRoomUserResult errcode:" + errcode);
}

function reqAddRoomAdmin(ridstr, idstr) {
  myapp.addroomadmin(ridstr, idstr, onAddRoomAdminResult);
}

function onAddRoomAdminResult(errcode) {
  console.info("onAddRoomAdminResult errcode:" + errcode);
}

function reqRemoveRoomAdmin(ridstr, idstr) {
  myapp.removeroomadmin(ridstr, idstr, onRemoveRoomAdminResult);
}

function onRemoveRoomAdminResult(errcode) {
  console.info("onRemoveRoomAdminResult errcode:" + errcode);
}

function reqSendRoomMessage(ridstr, message) {
  myapp.sendroommessage(ridstr, message, onSendRoomMessageResult);
}

function onSendRoomMessageResult(errcode) {
  console.info("onSendRoomMessageResult errcode:" + errcode);
}

function reqRoomList() {
  myapp.reqroomlist(onRoomListResult);
}

function onRoomListResult(errcode, data) {
  console.info("onRoomListResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(data));
  clearRoomList();
  roomdata = {};
  for(var i = 0; i < data.length; i ++) {
    //console.info("onRoomListResult data:" + JSON.stringify(data[i]));
    addRoom(data[i]);
    roomdata[data[i].rid] = data[i];

    reqRoomPresenceList(data[i].rid);
  }
}

function reqRoomPresenceList(ridstr) {
  myapp.reqroompresencelist(ridstr, onRoomPresenceListResult);
}

function onRoomPresenceListResult(errcode, data) {
  console.info("onRoomPresenceListResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(data));
  if(errcode == 0)
    for(var i = 0; i < data.length; i ++) {
      //console.info("onRoomPresenceListResult data:" + JSON.stringify(data[i]));
      addPresence(data[i]);
    }
}

function onRoomPresence(jsondata) {
  console.info("onRoomPresence jsondata:" + JSON.stringify(jsondata));
  addPresence(jsondata);
  if(jsondata.presencetype == PresenceType.PresenceType_Subscribed){
    reqRoomList();
  }
}

function onRoomMessage(jsondata) {
  console.info("onRoomMessage jsondata:" + JSON.stringify(jsondata));
  // var msgarray = messagelist[jsondata.who];
  // if(msgarray == null){
  //   msgarray = new Array();
  //   messagelist[jsondata.who] = msgarray;
  // }
  // msgarray[msgarray.length] = jsondata;
  if(jsondata.who != userdata.id)
  addRoomMessage(jsondata);
  else if(jsondata.platform != myapp.platform)
  addRoomSendMessage(jsondata);
}

function reqRoomUserList(strrid) {
  myapp.reqroomuserlist(strrid, onRoomUserList);
}

function onRoomUserList(errcode, jsondata) {
  console.info("onRoomUserList errcode:" + errcode);
  if(errcode == 0){
    console.info(JSON.stringify(jsondata));
    for(var i = 0; i < jsondata.length; i ++) {
      //console.info("onRoomUserList data:" + JSON.stringify(jsondata[i]));
      addRoomUser(jsondata[i]);
      var tmp = roomuserdata[jsondata[i].rid];
      if(tmp == null)
        tmp = {};
      tmp[jsondata[i].dataid] = jsondata[i];
      roomuserdata[jsondata[i].rid] = tmp;
    }
  }
}
//room end

//search start
function reqSearchUserById(idstr) {
  myapp.reqsearchuserbyid(idstr, onSearchUserByIdResult);
}

function onSearchUserByIdResult(errcode, data) {
  console.info("onSearchUserByIdResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(data));
  if(errcode == 0 && data) {
    clearSearchUserById();
    // for(var i = 0; i < data.length; i ++) {
    //   console.info("onSearchUserByIdResult data:" + JSON.stringify(data[i]));
    //   addSearchUserById(data[i]);
    // }
    addSearchUserById(data);
  }
}

function reqSearchUserByNickname(nickname) {
  myapp.reqsearchuserbynickname(nickname, onSearchUserByNicknameResult);
}

function onSearchUserByNicknameResult(errcode, data) {
  console.info("onSearchUserByNicknameResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(data));
  clearSearchUserByNickname();
  for(var i = 0; i < data.length; i ++) {
    //console.info("onSearchUserByNicknameResult data:" + JSON.stringify(data[i]));
    addSearchUserByNickname(data[i]);
  }
}

function reqSearchRoom(idstr) {
  myapp.reqsearchroom(idstr, onSearchRoomResult);
}

function onSearchRoomResult(errcode, data) {
  console.info("onSearchRoomResult errcode:" + errcode);
  console.info("data:" + JSON.stringify(data));
  clearSearchRoom();
  for(var i = 0; i < data.length; i ++) {
    //console.info("onSearchRoomResult data:" + JSON.stringify(data[i]));
    addSearchRoom(data[i]);
  }
}
//search end
