var markers = [];
var map;


function urlParam(name) {
    var url = window.location.href;
    var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(url);
    if (!results) {
        return undefined;
    }
    return decodeURIComponent(results[1]) || undefined;
}

function buildHTML(location) {
  if(location === undefined) {
    return "Unknown location"
  }

  body = ""
  for(item in location) {
    body += `<li><strong>${item}:</strong> <code>${location[item]}</code></li>`
  }

  return `<ul>${body}</ul>`
}

function updateLocationCount(count) {
  var el = document.getElementById('location-count')
  var suffix = "location"
  if(count != 1) {
    suffix += "s"
  }
  el.innerText = `${count} ${suffix}`
}

function buildElementForLocation(location) {
  var el = document.createElement('div');
  el.className = 'marker';
  if("locationType" in location) {
    if(location.locationType.includes("store")) {
      el.classList.add("store")
    } else if(location.locationType.includes("supercharger")) {
      el.classList.add("supercharger")
    } else if(location.locationType.includes("service")) {
      el.classList.add("service")
    } else if(location.locationType.includes("standard charger") || location.locationType.includes("destination charger")) {
      el.classList.add("charger")
    }
  }
  if("openSoon" in location && location.openSoon) {
    el.classList.add("coming-soon")
  }
  return el
}

function query(event) {
  if(event) {
    event.preventDefault()
  }

  Pace.start()
  fetch("/graphql", {
    method: "POST",
    body: this.query.value,
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/graphql"
    }
  }).then(function(response) {
    if(response.ok) {
      return response.json().then(function(json) {
        // When we get a new result set, clear the markers
        markers.forEach(function(marker) {
          marker.remove()
        })

        var coordinates = json.data.locations.edges.map(function(edge) {
          return [edge.node.longitude, edge.node.latitude]
        });

        markers = []
        json.data.locations.edges.forEach(function(edge) {
          var el = buildElementForLocation(edge.node)
          var popup = new mapboxgl.Popup({
            offset: [0, -36]
          }).setHTML(buildHTML(edge.node))

          var marker = new mapboxgl.Marker(el, {
            offset: [-12, -33]
          }).setLngLat([edge.node.longitude, edge.node.latitude])
          .setPopup(popup)
          .addTo(map)

          console.log(marker)

          markers.push(marker)
        })

        var bounds;

        if(coordinates.length == 1) {
          bounds = new mapboxgl.LngLatBounds(coordinates[0], coordinates[0])
        } else if(coordinates.length > 1) {
          bounds = coordinates.reduce(function(bounds, coord) {
            return bounds.extend(coord);
          }, new mapboxgl.LngLatBounds(coordinates[0], coordinates[0]));
        }

        map.fitBounds(bounds, {
          padding: 100,
          linear: false,
          maxZoom: 3
        });

        updateLocationCount(markers.length)
        Pace.stop()
      })
    } else {
      Pace.stop()
      console.log('Network response was not ok.');
    }
  })
  .catch(function(error) {
    console.log('There has been a problem with your fetch operation: ' + error.message);
    Pace.stop()
  });
}

Pace.on('start', function() {
  document.getElementById('header').classList.add("loading")
})
Pace.on('hide', function() {
  document.getElementById('header').classList.remove("loading")
})

window.onload = function() {
  map = new mapboxgl.Map({
      container: 'map',
      style: 'mapbox://styles/mapbox/streets-v9',
      zoom: 3,

      // Center of the US
      center: [-98.585522, 39.8333333]
  });
  var form = document.getElementsByTagName('form')[0];
  var queryParam = urlParam("query")
  if(queryParam != undefined) {
    form.query.value = queryParam
  }

  form.addEventListener('submit', query)

  form.query.addEventListener('keydown', function(e) {
  	if(e.keyCode == 13 && e.metaKey) {
  		query.apply(this.form)
  	}
  })

  // Fetch initial placeholder query data
  query.apply(form)

  window.addEventListener('keyup', function(e) {
    // s and / key
  	if(e.keyCode == 83 || e.keyCode == 191) {
  		form.query.focus()
  	}
  })
}
