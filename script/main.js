$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	$("#nocm").focus(function(){
		var value = $("#nocm").val();
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
			//nocm = "";
		} else {
			$("#datapasien").html("No. CM tidak lengkap");
			//nocm = "";
		} 
	});
	
	$("#form1").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
			token = localStorage.getItem("token")
			$.ajax({
				type: 'post',
				url: '/getcm',
				data: "nocm="+nocm+"&token="+token,
				success:function(data){
				   $("#datapasien").html(data);
				   $("#nocm").prop("disabled", false);
				}
				
			})
		}
	});
	
	$("#form1").on("click", "#btnsub", function(){
		
		var nocm = $("#nocm").val();
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[type='radio'][name='ats']:checked").val();
		var iki = $("input[type='radio'][name='iki']:checked").val();
		var shift = $("input[type='radio'][name='shift']:checked").val();
		
		if (nocm == ""||namapts == ""||diag == ""||ats == ""||iki == ""||shift ==""){
			alert("Data Belum Lengkap");
		}else{
			$.ajax({
				type:'post',
				url:'inputdatapts',
				data:"nocm="+nocm+"&namapts="+namapts+"&diag="+diag+"&ats="+ats+"&iki="+iki+"&shift="+shift,
				success:function(){
					location.reload();
				}
				
			})
		}
		
	})
    
});
