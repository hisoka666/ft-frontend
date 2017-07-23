$(document).ready(function(){
	var nocm = $("#nocm").val()
	
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
	
	$("#navbar").on("click", "#btnsub", function(event){		
		event.preventDefault();
		
		var nocm = $("#nocm").val();
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[name='ats']:checked").val();
		var iki = $("input[name='iki']:checked").val();
		var shift = $("input[name='shift']:checked").val();
		var bagian = $("input[name='bagian']:checked").val();
		var baru = $("#baru").val();
		var dok = $("#email").val();

		switch (true) {
			case nocm === "":
			break;
			case namapts === "":
			break;
			case diag === "":
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Diagnosis belum diisi" +
		                        "</div>");
			// alert("Diagnosis harus diisi");
			break;
			case ats == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"ATS belum diisi" +
		                        "</div>");
			// alert("ATS harus diisi");
			break;
			case iki == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Golongan IKI belum diisi" +
		                        "</div>");
			// alert("Golongan IKI harus diisi");
			break;
			case shift == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Shift belum diisi" +
		                        "</div>");
			// alert("Shift harus diisi");
			break;
			case bagian == null:
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Bagian belum diisi" +
		                        "</div>");
			// alert("Bagian harus diisi");
			break;
			default:
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
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script);
				}else{
					$("tbody").prepend(js.script);
					refreshNumber();
					$("#nocm").val('');
					$("#datapasien").html('');
					$("tbody tr:eq(101)").remove();
					popModalWarning("Sukses", "Berhasil menambahkan data", "")
					
				}
			});
		}
	});

	$("#navbar").on("click", "#editbut", function(event){
		event.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		token = localStorage.getItem("token");
			console.log("before show mymodal works");

			$.post("editentri",{
				token: token,
				link: link
			},
			function(data){
				var js = JSON.parse(data);
				$("#mymodal").html(js.script);
				var ats = 'input[name="modats"][value=' + js.content.ats + ']';
				var iki = 'input[name="modiki"][value=' + js.content.iki + ']';
				var shift = 'input[name="modshift"][value=' + js.content.shift + ']';
				var bagian = 'input[name="modbagian"][value=' + js.content.bagian + ']';
			
				$('input[name="entri"]').val(js.content.link);
				$('input[name="namapasien"]').val(js.content.nama);
				$('input[name="diagnosis"]').val(js.content.diag);
				$('input[name="urutan"]').val(indexrow)
			
				if (js.content.bagian == ""){
 		    		$(ats).prop('checked', true);
   					$(iki).prop('checked', true);
   					$(shift).prop('checked', true);	
			
				}else{
 		    		$(ats).prop('checked', true);
   					$(iki).prop('checked', true);
   					$(shift).prop('checked', true);
   					$(bagian).prop('checked', true);

				};
				$("#mymodal").modal();
			});
	});

$('body').on('hide.bs.modal', "#mymodal" , function(){
	console.log("Removal mymodal works");
	$(this).removeData('bs.modal').find(".modal-content").empty();
	});
$('body').on('hide.bs.modal', "#modwarning" , function(){
	console.log("removal modwarning works");
	$(this).removeData('bs.modal');
	});

$("body").on("click", "#confirmedit", function(e){
		e.preventDefault();
		
		var link = $("#modentri").val();
		var namapts = $('input[name="namapasien"]').val();
		var diag = $('input[name="diagnosis"]').val();
		var ats = $("input[name='modats']:checked").val();
		var bagian = $("input[name='modbagian']:checked").val();
		var iki = $("input[name='modiki']:checked").val();
		var shift = $("input[name='modshift']:checked").val();
		var indexrow = $('input[name="urutan"]').val();
		var urutan = "tbody tr:eq(" + $('input[name="urutan"]').val() + ")";

		
		switch (true) {
			case namapts === "":
			popModalWarning();
			break;
			case diag === "":
			popModalWarning();
			break;
			case ats == null :
			break;
			case iki == null :
			break;
			case shift == null :
			break;
			case bagian == null:
			break;
			default:
			$.post("confedit",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				link: link
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					popModalWarning("Kesalahan Pada Server", "Telah terjadi kesalahan pada server. Mohon ulangi proses sebelumnya");
				}else{
					$(urutan).replaceWith(js.script);
					console.log("Editing jalan")
					refreshNumber();
					$("#mymodal").modal('hide');
					popModalWarning("Edit Entri", "Berhasil mengubah entri", "");
					
				}
			});
		}

	});

var refreshNumber = function(){
	
$("tr").find(".nourut").each(function(index, elem){
		num = index + 1;
		$(elem).html(num)
	})
}

var popModalWarning = function(title, msg, prop){
		$(".modal-title").html(title);
		$("#message").html(msg);
		if (prop == ""){
			$("#extrabut").hide();
			$("#modwarning").modal();
		} else {
			$("#extrabut").html(prop);
			$("#extrabut").show();
			$("#modwarning").modal();
		}			
};

$("#navbar").on("click", "#delbut", function(e){
		e.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		var urutan = "tbody tr:eq(" + indexrow + ")";
		var token = localStorage.getItem("token");
		popModalWarning("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus");
		$("body").one("click", "#extrabut", function(){
			
			$.post("confdel", {
				link: link,
				token: token
			}, function(data){
				var js = JSON.parse(data);
				if (js.token == "OK"){
					$(urutan).remove();
					console.log("Removed : "+ urutan);
					refreshNumber();
					$("#modwarning").modal('hide');
				}else{
					alert("Terjadi kesalahan");
				}
			});
		});
	});






$("#navbar").on("click", "#homebutton", function(e){
		token = localStorage.getItem("token");
		email = $("#email").val();
		e.preventDefault();
		$.post("/firstentries", {
			token: token,
			email: email
		},function(data){
			var js = JSON.parse(data);
			// console.log(js.script);
			if (js.token == "OK") {
				$("#tabelutama").html(js.script);
				// removeModal("#modwarning")
			}else{
				popModalWarning("Peringatan", "Terjadi kesalahan pada server. Hubungi admin");
			}
			

		})

	});

$("#navbar").on("click", "#tglbut", function(e){
	token = localStorage.getItem("token");
	var link = $(this).offsetParent().children().first().html();
	e.preventDefault();
	console.log(link);

	$.post("edittgl", {
		token: token,
		link: link
	}, function(data){
		var js = JSON.parse(data)
		$("#mymodal").html(js.script)
		$("#mymodal").modal()
		console.log(js.script)

	})
});

});


