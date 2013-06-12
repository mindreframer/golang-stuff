
Node = function(g, data) {
  var that = new Object();
	that.data = data;
	that.gob_addr = data.Addr;
	var m = /^(.*):(.*)$/.exec(that.gob_addr);
	that.json_addr = m[1] + ":" + (1 + parseInt(m[2]));
	that.hexpos = "";
	_.each($.base64.decode2b(data.Pos), function(b) {
	  var thishex = b.toString(16);
		while (thishex.length < 2) {
		  thishex = "0" + thishex;
		}
	  that.hexpos += thishex;
	});
	var x_y = g.getPos(data.Pos);
	that.x = x_y[0];
	that.y = x_y[1];
	while (that.hexpos.length < 32) {
	  that.hexpos = "0" + that.hexpos;
	}
	return that;
}

God = function() {
  var that = new Object();
	that.maxPos = Big(1).times(Big(256).pow(16));
	that.api_endpoint_templ = _.template($("#api_endpoint_templ").text());
	that.api_endpoint_item_templ = _.template($("#api_endpoint_item_templ").text());
	that.result_templ = _.template($("#result_templ").text());
	that.node_link_templ = _.template($("#node_link_templ").text());
	that.cx = 1000;
	that.cy = 1000;
	that.r = 800;
	that.api_endpoints = {{.ApiMethods}};
	that.last_route_redraw = new Date().getTime();
	that.last_route_update = new Date().getTime();
	that.last_meta_redraw = new Date().getTime();
	that.last_meta_update = new Date().getTime();
	that.getPos = function(b64) {
		var angle = Big(3 * Math.PI / 2).plus($.base64.decode2big(b64).div(that.maxPos).times(Math.PI * 2)).toFixed();
		return [that.cx + Math.cos(angle) * that.r, that.cy + Math.sin(angle) * that.r];
	};
  that.drawChord = function() {
		var stage = new createjs.Stage(document.getElementById("chord"));

		var circle = new createjs.Shape();
		circle.graphics.beginStroke(createjs.Graphics.getRGB(0,0,0)).drawCircle(that.cx, that.cy, that.r);
		stage.addChild(circle);

    var dash = new createjs.Shape();
		dash.graphics.beginStroke(createjs.Graphics.getRGB(0,0,0)).moveTo(that.cx, that.cy - that.r - 30).lineTo(that.cx, that.cy - that.r + 30);
		stage.addChild(dash);

		if (that.last_route_update > that.last_route_redraw) {
			$("#nodes .node").remove();
		}
		for (var addr in that.node_by_addr) {
			var route = that.node_by_addr[addr];    
			var spot = new createjs.Shape();
			spot.graphics.beginStroke(createjs.Graphics.getRGB(0,0,0)).beginFill(createjs.Graphics.getRGB(0,0,0)).drawCircle(route.x, route.y, 20);
			stage.addChild(spot);
			var label = new createjs.Text(route.hexpos + "@" + route.gob_addr, "bold 25px Courier");
			label.x = route.x + 30;
			label.y = route.y - 10;
			stage.addChild(label);
			if (that.last_route_update > that.last_route_redraw) {
				$("#nodes").append(that.node_link_templ({ node: route }));
			}
		}
    if (that.last_route_update > that.last_route_redraw) {
			$(".node").click(function(e) {
			  that.selectNodeWithAddr($(e.target).parent().attr('data-addr'));
			});
			that.last_route_redraw = new Date().getTime();
		}

    var fade = 300;
    var newAnimations = [];
		var now = new Date().getTime()
		_.each(that.animations, function(ani) {
		  if (ani.ttl > now) {
			  newAnimations.push(ani);
			}
			var age = now - ani.created;
			var left = ani.ttl - now;
			var alpha = 1;
			var len = 1;
			var sx = ani.source[0];
			var sy = ani.source[1];
			var dx = ani.destination[0];
			var dy = ani.destination[1];
			var len;
			if (left < fade) {
			  alpha = left / fade;
				if (ani.key == null) {
					len = left / fade;
					sx = ani.destination[0] - (ani.destination[0] - ani.source[0]) * len;
					sy = ani.destination[1] - (ani.destination[1] - ani.source[1]) * len;
				}
			}
			if (age < fade && ani.key == null) {
			  len = age / fade;
				dx = ani.source[0] + (ani.destination[0] - ani.source[0]) * len;
				dy = ani.source[1] + (ani.destination[1] - ani.source[1]) * len;
			}
		  var line = new createjs.Shape()
			var gr = line.graphics.beginStroke(createjs.Graphics.getRGB(ani.color[0], ani.color[1], ani.color[2], alpha)).setStrokeStyle(ani.strokeWidth, ani.caps);
			if (ani.key != null) {
				gr.moveTo(sx, sy).quadraticCurveTo(ani.key[0], ani.key[1], dx, dy);
			} else {
				gr.moveTo(sx, sy).lineTo(dx, dy);
			}
			stage.addChild(line);
		});
		that.animations = newAnimations;

		stage.update();

    if (that.last_meta_update > that.last_meta_redraw) {
		  if (that.node != null) {
				$("#node_json_addr").text(that.node.json_addr);
				$("#node_gob_addr").text(that.node.gob_addr);
				$("#node_pos").text(that.node.hexpos);
				$("#node_owned_keys").text(that.node.data.OwnedEntries);
				$("#node_held_keys").text(that.node.data.HeldEntries);
				$("#node_load").text(that.node.data.Load);
			}
			that.last_meta_redraw = new Date().getTime();
		}
	};
	that.selectNodeWithAddr = function(addr) {
	  for (var a in that.node_by_addr) {
		  if (a == addr) {
			  that.node = that.node_by_addr[a];
				that.last_meta_update = new Date().getTime();
			  $("#node_container").css('display', 'block');
			  $("#hide_node_container").click(function(e) {
				  e.preventDefault();
				  $("#node_container").css('display', 'none');
				});
			}
		}
	};
	that.animations = [];
	that.animate = function(item) {
	  if (that.animations.length < 100) {
		  that.animations.push(item);
		}
	};
	that.opened_sockets = {};
	that.node_by_addr = {};
	that.node = null;
	that.open_socket = function(addr) {
	  addr = addr.replace('localhost', '127.0.0.1');
	  if (that.opened_sockets[addr] == null) {
			that.opened_sockets[addr] = true;
			$.websocket("ws://" + addr + "/ws", {
				open: function() { 
					console.log("socket to " + addr + " opened");
				},
				close: function() { 
				  delete(that.opened_sockets[addr]);
					delete(that.node_by_addr[addr]);
					console.log("socket to " + addr + " closed");
				},
				events: {
					RingChange: function(e) {
						_.each(e.data.routes, function(r) {
							var node = new Node(that, r);
							if (that.node_by_addr[node.json_addr] == null) {
								that.open_socket(node.json_addr);
							}
						});
						var newNode = new Node(that, e.data.description);
						that.node_by_addr[newNode.json_addr] = newNode;
						if (that.node != null && that.node.json_addr == newNode.json_addr) {
							that.node = newNode;
						}
						that.last_route_update = new Date().getTime();
						that.last_meta_update = new Date().getTime();
					},
					Comm: function(e) {
						var item = {
							source: that.getPos(e.data.source.Pos),
							destination: that.getPos(e.data.destination.Pos),
							ttl: new Date().getTime() + 400,
							color: [0,0,200],
							strokeWidth: 3,
							caps: 0,
							created: new Date().getTime(),
						};
						if (/Notify/.exec(e.data.type) != null) {
							item.color = [200,200,0];
						}
						if (/Ping/.exec(e.data.type) != null) {
							item.color = [200,0,200];
						}
						if (/HashTree/.exec(e.data.type) != null) {
							item.color = [200,0,0];
						}
						if (e.data.key != null) {
							item.key = that.getPos(e.data.key);
						}
						if (e.data.sub_key != null) {
							item.sub_key = that.getPos(e.data.sub_key);
						}
						that.animate(item);
					},
					Sync: function(e) {
						var item = {
							source: that.getPos(e.data.source.Pos),
							destination: that.getPos(e.data.destination.Pos),
							ttl: new Date().getTime() + 400,
							color: [150,0,150],
							strokeWidth: 5,
							created: new Date().getTime(),
							caps: 1,
						};
						that.animate(item);
					},
					Clean: function(e) {
						var item = {
							source: that.getPos(e.data.source.Pos),
							destination: that.getPos(e.data.destination.Pos),
							ttl: new Date().getTime() + 400,
							created: new Date().getTime(),
							color: [50,150,0],
							strokeWidth: 4,
							caps: 2,
						};
						that.animate(item);
					},
				},
			});
		}
	};
	that.generate_example_param = function(param) {
	  if (param == 'bool') {
		  return 'false';
		} else if (param == '[]byte') {
		  return '$.base64.encode("some bytes")';
		} else if (param == 'int') {
		  return '0';
		} else {
		  return '"Unknown parameter type ' + param + ', please implement example code for it!"';
		}
	};
	that.generate_example_params = function(endp) {
	  var rval = '{\n';
		for (var key in endp.parameter) {
		  rval += '    ' + key + ': ' + that.generate_example_param(endp.parameter[key]) + ',\n';
		}
		rval += '  }';
		return rval;
	};
	that.display_endpoint_form = function(endp) {
		$("#code_container").html(that.api_endpoint_templ({ api_endpoint: endp }));
    var generated_code = "$.ajax('/rpc/DHash." + endp.name + "', {\n  type: 'POST',\n  contentType: 'application/json; charset=UTF-8',\n  data: JSON.stringify(" + that.generate_example_params(endp) + "),\n  success: function(data) {\n    displayResult(data);\n  },\n  dataType: 'json',\n});";
		$("#code").val(generated_code);
		$("#execute").click(function(ev) {
		  eval($("#code").val());
		});
	};
	that.start = function() {
		window.setInterval(that.drawChord, 40);
		that.open_socket(document.location.hostname + ":" + document.location.port);
		_.each(that.api_endpoints, function(api_endpoint) {
		  $("#endpoints").append(that.api_endpoint_item_templ({ api_endpoint: api_endpoint }));
		});
		$("#endpoints li").click(function(ev) {
		  that.display_endpoint_form(that.api_endpoints[$(ev.target).attr("data-endpoint-name")]);
		});
	};
	return that;
};

g = new God();

result = null;

function b64decode(obj) {
  if (obj instanceof Array) {
	  var rval = [];
		for (var i = 0; i < obj.length; i++) {
			rval.push(b64decode(obj[i]));
		}
		return rval;
	} else {
		var rval = {};
		for (var name in obj) {
			var value = obj[name];
			if (typeof(value) == "string") {
			  try {
					rval[name] = $.base64.decode2s(value);
				} catch (e) {
				  rval[name] = value;
				}
			} else if (value != null && typeof(value) == "object") {
				rval[name] = b64decode(value);
			} else {
				rval[name] = value;
			}
		}
		return rval;
	}
};

function displayResult(data) {
  result = data;
  $("#result_container").html(g.result_templ({ data: result }));
	$("#decode").click(function(ev) {
		displayResult(b64decode(result));
		$("#decode").remove();
	});
};

$(function() {
  g.start();
});

