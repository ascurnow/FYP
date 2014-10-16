window.onload = function(){
	cookieValue = getCookie('cookie:username');
	document.getElementById('ANAME0').innerHTML = cookieValue;
	document.getElementById('ANAME1').innerHTML = cookieValue;
}