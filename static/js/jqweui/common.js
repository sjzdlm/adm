function gourl(url) {
	$.showLoading("正在处理...");
	window.location.href = url;
}

// / <summary>
// / 判断是否是日期
// / </summary>
// / <param name="sDate">日期字符串</param>
// / <returns>返回是否(bool)</returns>
function IsDate(sDate) {
	var sRegex = /^(\d{4})-(\d{2})-(\d{2})$/;
	var bResult = sDate.match(reg);
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / <summary>
// / 判断字符串是否为空
// / </summary>
// / <param name="sNullOrEmpty">空字符串</param>
// / <returns>返回是否(bool)</returns>
function IsNullEmpty(sNullOrEmpty) {
	if (sNullOrEmpty.length == '' || sNullOrEmpty.length <= 0) {
		return false;
	} else {
		return true;
	}
}
function IsCurrent(sCurrent) {
	var bResult1 = sCurrent.match("[^0-9.-]");
	var bResult2 = sCurrent.match("[[0-9]*[.][0-9]*[.][0-9]*");
	var bResult3 = sCurrent.match("[[0-9]*[-][0-9]*[-][0-9]");
	var bResult4 = sCurrent
			.match("(^([-]|[.]|[-.]|[0-9])[0-9]*[.]*[0-9]+$)|(^([-]|[0-9])[0-9]*$)");
	if (bResult1 != null || bResult2 != null || bResult3 != null
			|| bResult4 == null) {
		return false;
	} else {
		return true;
	}
}
// / <summary>
// / 判断是否是数字
// / </summary>
// / <param name="sNum">数字字符串</param>
// / <returns>返回是否(bool)</returns>
function IsNumeric(sNum) {
	var bResult = sNum.match("^(-|\\+)?\\d+(\\.\\d+)?$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// 正整数
function isPInt(str) {
	var g = /^[1-9]*[1-9][0-9]*$/;
	return g.test(str);
}
function isInteger(str) {
	if (/^-?\d+$/.test(str)) {
		return true;
	}
	return false;
}

function isFloat(str) {
	if (/^(-?\d+)(\.\d+)?$/.test(str)) {
		return true;
	}
	return false;
}

// / <summary>
// / 判断是否是URL
// / </summary>
// / <param name="sUrl">URL字符串</param>
// / <returns>返回是否(bool)</returns>
function IsUrl(sUrl) {
	var bResult = sUrl
			.match("http(s)?://([\\w-]+\\.)+[\\w-]+(/[\\w- ./?%&=]*)?");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是否是MAIL
// / </summary>
// / <param name="sMail">MAIL字符串</param>
// / <returns>返回是否(bool)</returns>
function IsMail(sMail) {
	var bResult = sMail
			.match("\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是否是邮编
// / </summary>
// / <param name="sPostCode">邮编字符串</param>
// / <returns>返回是否(bool)</returns>
function IsPostCode(sPostCode) {
	var bResult = sPostCode.match("^\\d{6}$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是否是电话号码
// / </summary>
// / <param name="sTelephone">电话号码字符串</param>
// / <returns>返回是否(bool)</returns>
function IsTelephone(sTelephone) {
	var bResult = sTelephone.match("^(\\(\\d{3}\\)|\\d{3}-)?\\d{8}$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是否是手机号码
// / </summary>
// / <param name="sMobile">手机号码字符串</param>
// / <returns>返回是否(bool)</returns>
function IsMobile(sMobile) {
	var bResult = sMobile.match("^\\d{11}$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是否身份证
// / </summary>
// / <param name="sSimNum">数字字符串</param>
// / <returns>返回是否(bool)</returns>
function IsIDCard(sIDCard) {
	var bResult = sIDCard.match("^\\d{15}|\\d{18}$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
// / 判断是中英表达式
// / </summary>
// / <param name="sCE">中英文表达式字符串</param>
// / <returns>返回是否(bool)</returns>
function IsCE(sCE) {
	var bResult = sCE.match("^[a-zA-Z\\u4E00-\\u9FA5\\uF900-\\uFA2D]+$");
	if (bResult == null) {
		return false;
	} else {
		return true;
	}
}
function isChinese(str) {
	var str = str.replace(/(^\s*)|(\s*$)/g, '');
	if (!(/^[\u4E00-\uFA29]*$/.test(str) && (!/^[\uE7C7-\uE7F3]*$/.test(str)))) {
		return false;
	}
	return true;
}
function isImg(str) {
	var objReg = new RegExp("[.]+(jpg|jpeg|swf|gif)$", "gi");
	if (objReg.test(str)) {
		return true;
	}
	return false;
}
