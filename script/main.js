$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	// $("#nocm").focus(function(){
	// 	var value = $("#nocm").val();
	// 	if (value == ""){
	// 		$("#datapasien").html("Masukkan No. CM");
	// 		//nocm = "";
	// 	} else {
	// 		$("#datapasien").html("No. CM tidak lengkap");
	// 		//nocm = "";
	// 	} 
	// });
	
	$("#navbar").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
            token = localStorage.getItem("token");

			$.post("getcm", {
				nocm: nocm,
				token: token
			},
			function(data){
				var js = JSON.parse(data)

				if (js.token == "OK"){
					$("#datapasien").html(js.script);
					$("#nocm").prop("disabled", false);
				} else {
					alert("Terjadi kesalahan: " + js.token)
					$("#nocm").prop("disabled", false);
				}			    

			});
        }

	});
	
	$("#navbar").on("click", "#btnsub", function(){
		
		var nocm = $("#nocm").val();
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[type='radio'][name='ats']:checked").val();
		var iki = $("input[type='radio'][name='iki']:checked").val();
		var shift = $("input[type='radio'][name='shift']:checked").val();
		var bagian = $("input[type='radio'][name='bagian']:checked").val();
		var baru = $("#baru").val();
		// var auth2 = gapi.auth2.getAuthInstance();
		var profile = gapi.auth2.getBasicProfile();
		var dok = profile.getEmail();
		if (nocm == ""||namapts == ""||diag == ""||ats == ""||iki == ""||shift ==""||bagian==""){
			alert("Data Belum Lengkap");
		}else{
			$.post("inputdata",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				dok: dok,
				baru: baru
			},
			function(data){
				$("tbody").prepend(data)
			})
			// $.post({
			// 	type:'post',
			// 	url:'inputdata',
			// 	data:"nocm="+nocm+"&namapts="+namapts+"&diag="+diag+"&ats="+ats+"&iki="+iki+"&shift="+shift,
			// 	success:function(){
			// 		location.reload();
			// 	}
				
			// })
		}
		
	})

})