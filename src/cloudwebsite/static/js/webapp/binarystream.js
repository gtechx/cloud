
function stringToBytes(str) {
  return new TextEncoder("utf-8").encode(str);
}

function bytesToString(buffer) {
  return new TextEncoder("utf-8").decode(buffer);
}

Date.prototype.format = function (fmt) { // author: meizz
  if(fmt == "" || fmt == undefined || fmt == null)
      fmt = "yyyy-MM-ddThh:mm:ss.SSS+th:tm";
  var date = {
      "M+": this.getMonth() + 1, // 月份
      "d+": this.getDate(), // 日
      "h+": this.getHours(), // 小时
      "m+": this.getMinutes(), // 分
      "s+": this.getSeconds(), // 秒
      "q+": Math.floor((this.getMonth() + 3) / 3), // 季度
      "S+": this.getMilliseconds(), // 毫秒
      "th+": this.getTimezoneOffset() / -60, //时区
		  "tm+": 0 //时区的分
  };
  if (/(y+)/i.test(fmt))
      fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
  for (var k in date)
      if (new RegExp("(" + k + ")").test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? date[k] : ("000" + date[k]).substr(-RegExp.$1.length));
      }
      return fmt;
}

var BinaryStream = {
  new: function () {
    var stream = {};
    stream.length = 0;
    stream.cur = 0;
    
    if(arguments[0] && arguments[0] instanceof ArrayBuffer){
      stream.buffer = arguments[0];
      stream.length = arguments[0].byteLength;
      stream.cap = arguments[0].byteLength;
    } else {
      stream.cap = 512;
      stream.buffer = new ArrayBuffer(stream.cap);
    }
    stream.dataview = new DataView(stream.buffer);
    stream.int8view = new Int8Array(stream.buffer);
    stream.uint8view = new Uint8Array(stream.buffer);
    stream.littleEndian = true;

    stream.getBuffer = function (){
      return stream.buffer.slice(0, stream.length);
    };

    stream.reset = function (){
      if(arguments[0] && arguments[0] instanceof ArrayBuffer){
        stream.buffer = arguments[0];
        stream.length = arguments[0].byteLength;
        stream.cap = arguments[0].byteLength;
        stream.dataview = new DataView(stream.buffer);
        stream.int8view = new Int8Array(stream.buffer);
        stream.uint8view = new Uint8Array(stream.buffer);
      } else {
        stream.length = 0;
      }
      stream.cur = 0;
      return stream;
    };

    stream.writeUint8 = function (data){
      stream.dataview.setUint8(stream.cur, data);
      stream.cur = stream.cur + 1;
      stream.length = stream.length + 1;
      return stream;
    };

    stream.writeInt8 = function (data){
      stream.dataview.setInt8(stream.cur, data);
      stream.cur = stream.cur + 1;
      stream.length = stream.length + 1;
      return stream;
    };

    stream.writeInt16 = function (data){
      stream.dataview.setInt16(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 2;
      stream.length = stream.length + 2;
      return stream;
    };

    stream.writeUint16 = function (data){
      stream.dataview.setUint16(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 2;
      stream.length = stream.length + 2;
      return stream;
    };

    stream.writeInt32 = function (data){
      stream.dataview.setInt32(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 4;
      stream.length = stream.length + 4;
      return stream;
    };

    stream.writeUint32 = function (data){
      stream.dataview.setUint32(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 4;
      stream.length = stream.length + 4;
      return stream;
    };

    stream.writeInt64 = function (data){
      if(data instanceof Long) {
        stream.int8view.set(data.toBytes(stream.littleEndian), stream.cur);
        stream.cur = stream.cur + 8;
        stream.length = stream.length + 8;
      }
      return stream;
    };

    stream.writeUint64 = function (data){
      if(data instanceof Long) {
        stream.int8view.set(data.toBytes(stream.littleEndian), stream.cur);
        stream.cur = stream.cur + 8;
        stream.length = stream.length + 8;
      }
      return stream;
    };

    stream.writeFloat32 = function (data){
      stream.dataview.setFloat32(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 4;
      stream.length = stream.length + 4;
      return stream;
    };

    stream.writeFloat64 = function (data){
      stream.dataview.setFloat64(stream.cur, data, stream.littleEndian);
      stream.cur = stream.cur + 8;
      stream.length = stream.length + 8;
      return stream;
    };

    stream.writeString = function (str){
      var strdata = new TextEncoder("utf-8").encode(str);
      stream.uint8view.set(strdata, stream.cur);
      stream.cur = stream.cur + strdata.byteLength;
      stream.length = stream.length + strdata.byteLength;
      return stream;
    };

    stream.writeArrayBuffer = function (data){
      if(data instanceof ArrayBuffer) {
        stream.uint8view.set(new Uint8Array(data), stream.cur);
        stream.cur = stream.cur + data.byteLength;
        stream.length = stream.length + data.byteLength;
      }
      return stream;
    };

    stream.writeArray = function (data){
      //if(data instanceof Array) {
        stream.uint8view.set(new Uint8Array(data), stream.cur);
        stream.cur = stream.cur + data.length;
        stream.length = stream.length + data.length;
      //}
      return stream;
    };

    //read
    stream.readUint8 = function (){
      var ret = stream.dataview.getUint8(stream.cur);
      stream.cur = stream.cur + 1;
      return ret;
    };

    stream.readInt8 = function (){
      var ret = stream.dataview.getInt8(stream.cur);
      stream.cur = stream.cur + 1;
      return ret;
    };

    stream.readUint16 = function (){
      var ret = stream.dataview.getUint16(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 2;
      return ret;
    };

    stream.readInt16 = function (){
      var ret = stream.dataview.getInt16(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 2;
      return ret;
    };

    stream.readUint32 = function (){
      var ret = stream.dataview.getUint32(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 4;
      return ret;
    };

    stream.readInt32 = function (){
      var ret = stream.dataview.getInt32(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 4;
      return ret;
    };

    stream.readUint64 = function (){
      // var arr = new Array(stream.uint8view.subarray(stream.cur, stream.cur + 8));
      // console.info("arr:"+arr.length);
      // console.info(arr);
      var ret = Long.fromBytes(stream.uint8view.subarray(stream.cur, stream.cur + 8), true, stream.littleEndian);//stream.dataview.getUint32(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 8;
      return ret;
    };

    stream.readInt64 = function (){
      var ret = Long.fromBytes(stream.uint8view.subarray(stream.cur, stream.cur + 8), false, stream.littleEndian);
      stream.cur = stream.cur + 8;
      return ret;
    };

    stream.readFloat32 = function (){
      var ret = stream.dataview.getFloat32(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 4;
      return ret;
    };

    stream.readFloat64 = function (){
      var ret = stream.dataview.getFloat642(stream.cur, stream.littleEndian);
      stream.cur = stream.cur + 8;
      return ret;
    };

    stream.readString = function (len){
      var ret = new TextDecoder("utf-8").decode(stream.buffer.slice(stream.cur, stream.cur + len));
      stream.cur = stream.cur + len;
      return ret;
    };

    stream.readStringAll = function (){
      var ret = new TextDecoder("utf-8").decode(stream.buffer.slice(stream.cur));
      stream.cur = stream.length - 1;
      return ret;
    };

    stream.readArrayBuffer = function (len){
      var ret = stream.buffer.slice(stream.cur, stream.cur + len);
      stream.cur = stream.cur + len;
      return ret;
    };

    return stream;
  }
}

//BinaryStream test
// var bs = BinaryStream.new();
// bs.writeInt8(0x01).writeInt16(0x0302).writeInt32(0x07060504).writeInt64(Long.fromString("0x0f0e0d0c0b0a0908", true, 16));
// console.info(bs.getBuffer());
// var newbs = BinaryStream.new(bs.getBuffer());

// console.info(newbs.readInt8().toString(16));
// console.info(newbs.readInt16().toString(16));
// console.info(newbs.readInt32().toString(16));
// console.info(newbs.readInt64().toString(16));

// var strbs = BinaryStream.new();
// strbs.writeString("abcdefg");
// console.info(strbs.getBuffer());
// console.info(strbs.cur);
// console.info(strbs.length);
