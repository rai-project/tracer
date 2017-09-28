package chrome

// from https://github.com/golang/go/blob/master/src/cmd/trace/trace.go

// See https://github.com/catapult-project/catapult/blob/master/tracing/docs/embedding-trace-viewer.md
// This is almost verbatim copy of:
// https://github.com/catapult-project/catapult/blob/master/tracing/bin/index.html
// on revision 623a005a3ffa9de13c4b92bc72290e7bcd1ca591.
const HTMLTemplate = `
<html>
<head>
<link href="/trace_viewer_html" rel="import">
<script>
(function() {
  var viewer;
  var url;
  var model;
  function load() {
    var req = new XMLHttpRequest();
    var is_binary = /[.]gz$/.test(url) || /[.]zip$/.test(url);
    req.overrideMimeType('text/plain; charset=x-user-defined');
    req.open('GET', url, true);
    if (is_binary)
      req.responseType = 'arraybuffer';
    req.onreadystatechange = function(event) {
      if (req.readyState !== 4)
        return;
      window.setTimeout(function() {
        if (req.status === 200)
          onResult(is_binary ? req.response : req.responseText);
        else
          onResultFail(req.status);
      }, 0);
    };
    req.send(null);
  }
  function onResultFail(err) {
    var overlay = new tr.ui.b.Overlay();
    overlay.textContent = err + ': ' + url + ' could not be loaded';
    overlay.title = 'Failed to fetch data';
    overlay.visible = true;
  }
  function onResult(result) {
    model = new tr.Model();
    var i = new tr.importer.Import(model);
    var p = i.importTracesWithProgressDialog([result]);
    p.then(onModelLoaded, onImportFail);
  }
  function onModelLoaded() {
    viewer.model = model;
    viewer.viewTitle = "trace";
  }
  function onImportFail() {
    var overlay = new tr.ui.b.Overlay();
    overlay.textContent = tr.b.normalizeException(err).message;
    overlay.title = 'Import error';
    overlay.visible = true;
  }
  document.addEventListener('DOMContentLoaded', function() {
    var container = document.createElement('track-view-container');
    container.id = 'track_view_container';
    viewer = document.createElement('tr-ui-timeline-view');
    viewer.track_view_container = container;
    viewer.appendChild(container);
    viewer.id = 'trace-viewer';
    viewer.globalMode = true;
    document.body.appendChild(viewer);
    url = '/jsontrace?{{PARAMS}}';
    load();
  });
}());
</script>
</head>
<body>
</body>
</html>
`
