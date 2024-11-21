"use strict";

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;

var _react = _interopRequireWildcard(require("react"));

var _core = require("@material-ui/core");

var _propTypes = _interopRequireDefault(require("prop-types"));

var _ThreeDEngine = _interopRequireDefault(require("./threeDEngine/ThreeDEngine"));

var _CameraControls = require("../camera-controls/CameraControls");

var _SelectionManager = require("./threeDEngine/SelectionManager");

var _reactResizeDetector = _interopRequireDefault(require("react-resize-detector"));

var _Recorder = require("./captureManager/Recorder");

var _Screenshoter = require("./captureManager/Screenshoter");

var _CaptureControls = require("../capture-controls/CaptureControls");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _getRequireWildcardCache(nodeInterop) { if (typeof WeakMap !== "function") return null; var cacheBabelInterop = new WeakMap(); var cacheNodeInterop = new WeakMap(); return (_getRequireWildcardCache = function _getRequireWildcardCache(nodeInterop) { return nodeInterop ? cacheNodeInterop : cacheBabelInterop; })(nodeInterop); }

function _interopRequireWildcard(obj, nodeInterop) { if (!nodeInterop && obj && obj.__esModule) { return obj; } if (obj === null || _typeof(obj) !== "object" && typeof obj !== "function") { return { "default": obj }; } var cache = _getRequireWildcardCache(nodeInterop); if (cache && cache.has(obj)) { return cache.get(obj); } var newObj = {}; var hasPropertyDescriptor = Object.defineProperty && Object.getOwnPropertyDescriptor; for (var key in obj) { if (key !== "default" && Object.prototype.hasOwnProperty.call(obj, key)) { var desc = hasPropertyDescriptor ? Object.getOwnPropertyDescriptor(obj, key) : null; if (desc && (desc.get || desc.set)) { Object.defineProperty(newObj, key, desc); } else { newObj[key] = obj[key]; } } } newObj["default"] = obj; if (cache) { cache.set(obj, newObj); } return newObj; }

function _extends() { _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; }; return _extends.apply(this, arguments); }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) { symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); } keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function asyncGeneratorStep(gen, resolve, reject, _next, _throw, key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { Promise.resolve(value).then(_next, _throw); } }

function _asyncToGenerator(fn) { return function () { var self = this, args = arguments; return new Promise(function (resolve, reject) { var gen = fn.apply(self, args); function _next(value) { asyncGeneratorStep(gen, resolve, reject, _next, _throw, "next", value); } function _throw(err) { asyncGeneratorStep(gen, resolve, reject, _next, _throw, "throw", err); } _next(undefined); }); }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function"); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, writable: true, configurable: true } }); if (superClass) _setPrototypeOf(subClass, superClass); }

function _setPrototypeOf(o, p) { _setPrototypeOf = Object.setPrototypeOf || function _setPrototypeOf(o, p) { o.__proto__ = p; return o; }; return _setPrototypeOf(o, p); }

function _createSuper(Derived) { var hasNativeReflectConstruct = _isNativeReflectConstruct(); return function _createSuperInternal() { var Super = _getPrototypeOf(Derived), result; if (hasNativeReflectConstruct) { var NewTarget = _getPrototypeOf(this).constructor; result = Reflect.construct(Super, arguments, NewTarget); } else { result = Super.apply(this, arguments); } return _possibleConstructorReturn(this, result); }; }

function _possibleConstructorReturn(self, call) { if (call && (_typeof(call) === "object" || typeof call === "function")) { return call; } return _assertThisInitialized(self); }

function _assertThisInitialized(self) { if (self === void 0) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return self; }

function _isNativeReflectConstruct() { if (typeof Reflect === "undefined" || !Reflect.construct) return false; if (Reflect.construct.sham) return false; if (typeof Proxy === "function") return true; try { Boolean.prototype.valueOf.call(Reflect.construct(Boolean, [], function () {})); return true; } catch (e) { return false; } }

function _getPrototypeOf(o) { _getPrototypeOf = Object.setPrototypeOf ? Object.getPrototypeOf : function _getPrototypeOf(o) { return o.__proto__ || Object.getPrototypeOf(o); }; return _getPrototypeOf(o); }

var styles = function styles() {
  return {
    container: {
      height: '100%',
      width: '100%'
    }
  };
};

var Canvas = /*#__PURE__*/function (_Component) {
  _inherits(Canvas, _Component);

  var _super = _createSuper(Canvas);

  function Canvas(props) {
    var _this;

    _classCallCheck(this, Canvas);

    _this = _super.call(this, props);
    _this.sceneRef = /*#__PURE__*/_react["default"].createRef();
    _this.cameraControls = /*#__PURE__*/_react["default"].createRef();
    _this.state = {
      modelReady: false
    };

    _this.constructorFromProps(props);

    _this.frameResizing = _this.frameResizing.bind(_assertThisInitialized(_this));
    _this.defaultCameraControlsHandler = _this.defaultCameraControlsHandler.bind(_assertThisInitialized(_this));
    _this.defaultCaptureControlsHandler = _this.defaultCaptureControlsHandler.bind(_assertThisInitialized(_this));
    return _this;
  }

  _createClass(Canvas, [{
    key: "constructorFromProps",
    value: function constructorFromProps(props) {
      if (props.captureOptions !== undefined) {
        this.captureControls = /*#__PURE__*/_react["default"].createRef();
      }
    }
  }, {
    key: "componentDidMount",
    value: function () {
      var _componentDidMount = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee() {
        var _this$props, data, cameraOptions, cameraHandler, captureOptions, backgroundColor, pickingEnabled, linesThreshold, hoverListeners, emptyHoverListener, setColorHandler, onMount, selectionStrategy, onSelection, updateStarted, updateEnded, dracoDecoderPath;

        return regeneratorRuntime.wrap(function _callee$(_context) {
          while (1) {
            switch (_context.prev = _context.next) {
              case 0:
                _this$props = this.props, data = _this$props.data, cameraOptions = _this$props.cameraOptions, cameraHandler = _this$props.cameraHandler, captureOptions = _this$props.captureOptions, backgroundColor = _this$props.backgroundColor, pickingEnabled = _this$props.pickingEnabled, linesThreshold = _this$props.linesThreshold, hoverListeners = _this$props.hoverListeners, emptyHoverListener = _this$props.emptyHoverListener, setColorHandler = _this$props.setColorHandler, onMount = _this$props.onMount, selectionStrategy = _this$props.selectionStrategy, onSelection = _this$props.onSelection, updateStarted = _this$props.updateStarted, updateEnded = _this$props.updateEnded, dracoDecoderPath = _this$props.dracoDecoderPath;
                this.threeDEngine = new _ThreeDEngine["default"](this.sceneRef.current, cameraOptions, cameraHandler, captureOptions, onSelection, backgroundColor, pickingEnabled, linesThreshold, hoverListeners, emptyHoverListener, setColorHandler, selectionStrategy, updateStarted, updateEnded, dracoDecoderPath);

                if (captureOptions) {
                  this.recorder = new _Recorder.Recorder(this.getCanvasElement(), captureOptions.recorderOptions);
                }

                _context.next = 5;
                return this.threeDEngine.start(data, cameraOptions, true);

              case 5:
                onMount(this.threeDEngine.scene);
                this.setState({
                  modelReady: true
                });
                this.threeDEngine.requestFrame();
                this.threeDEngine.setBackgroundColor(backgroundColor);

              case 9:
              case "end":
                return _context.stop();
            }
          }
        }, _callee, this);
      }));

      function componentDidMount() {
        return _componentDidMount.apply(this, arguments);
      }

      return componentDidMount;
    }()
  }, {
    key: "componentDidUpdate",
    value: function () {
      var _componentDidUpdate = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee2(prevProps, prevState, snapshot) {
        var _this$props2, data, cameraOptions, threeDObjects, backgroundColor;

        return regeneratorRuntime.wrap(function _callee2$(_context2) {
          while (1) {
            switch (_context2.prev = _context2.next) {
              case 0:
                if (!(prevProps !== this.props)) {
                  _context2.next = 8;
                  break;
                }

                _this$props2 = this.props, data = _this$props2.data, cameraOptions = _this$props2.cameraOptions, threeDObjects = _this$props2.threeDObjects, backgroundColor = _this$props2.backgroundColor;
                _context2.next = 4;
                return this.threeDEngine.update(data, cameraOptions, threeDObjects, this.shouldEngineTraverse(), backgroundColor);

              case 4:
                this.threeDEngine.requestFrame();
                this.setState({
                  modelReady: true
                });
                _context2.next = 9;
                break;

              case 8:
                this.setState({
                  modelReady: false
                });

              case 9:
              case "end":
                return _context2.stop();
            }
          }
        }, _callee2, this);
      }));

      function componentDidUpdate(_x, _x2, _x3) {
        return _componentDidUpdate.apply(this, arguments);
      }

      return componentDidUpdate;
    }()
  }, {
    key: "shouldComponentUpdate",
    value: function shouldComponentUpdate(nextProps, nextState, nextContext) {
      return nextState.modelReady || nextProps !== this.props;
    }
  }, {
    key: "componentWillUnmount",
    value: function componentWillUnmount() {
      this.threeDEngine.stop();
      this.sceneRef.current.removeChild(this.threeDEngine.getRenderer().domElement);
    }
  }, {
    key: "defaultCaptureControlsHandler",
    value: function defaultCaptureControlsHandler(action) {
      var captureOptions = this.props.captureOptions;

      if (this.recorder) {
        switch (action.type) {
          case _CaptureControls.captureControlsActions.START:
            this.recorder.startRecording();
            break;

          case _CaptureControls.captureControlsActions.STOP:
            var options = action.data.options;
            return this.recorder.stopRecording(options);

          case _CaptureControls.captureControlsActions.DOWNLOAD_VIDEO:
            {
              var _action$data = action.data,
                  filename = _action$data.filename,
                  _options = _action$data.options;
              return this.recorder.download(filename, _options);
            }
        }
      }

      if (captureOptions && captureOptions.screenshotOptions) {
        var _captureOptions$scree = captureOptions.screenshotOptions,
            quality = _captureOptions$scree.quality,
            pixelRatio = _captureOptions$scree.pixelRatio,
            resolution = _captureOptions$scree.resolution,
            filter = _captureOptions$scree.filter;

        switch (action.type) {
          case _CaptureControls.captureControlsActions.DOWNLOAD_SCREENSHOT:
            {
              var _filename = action.data.filename;
              (0, _Screenshoter.downloadScreenshot)(this.getCanvasElement(), quality, resolution, pixelRatio, filter, _filename);
              break;
            }
        }
      }
    }
  }, {
    key: "defaultCameraControlsHandler",
    value: function defaultCameraControlsHandler(action) {
      var defaultProps = {
        incrementPan: {
          x: 0.01,
          y: 0.01
        },
        incrementRotation: {
          x: 0.01,
          y: 0.01,
          z: 0.01
        },
        incrementZoom: 0.1,
        movieFilter: false
      };

      var mergedProps = _objectSpread(_objectSpread({}, defaultProps), this.props.cameraOptions.cameraControls);

      var incrementPan = mergedProps.incrementPan,
          incrementRotation = mergedProps.incrementRotation,
          incrementZoom = mergedProps.incrementZoom,
          movieFilter = mergedProps.movieFilter;

      if (this.threeDEngine) {
        switch (action) {
          case _CameraControls.cameraControlsActions.PAN_LEFT:
            this.threeDEngine.cameraManager.incrementCameraPan(-incrementPan.x, 0);
            break;

          case _CameraControls.cameraControlsActions.PAN_RIGHT:
            this.threeDEngine.cameraManager.incrementCameraPan(incrementPan.x, 0);
            break;

          case _CameraControls.cameraControlsActions.PAN_UP:
            this.threeDEngine.cameraManager.incrementCameraPan(0, -incrementPan.y);
            break;

          case _CameraControls.cameraControlsActions.PAN_DOWN:
            this.threeDEngine.cameraManager.incrementCameraPan(0, incrementPan.y);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_UP:
            this.threeDEngine.cameraManager.incrementCameraRotate(0, incrementRotation.y, undefined);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_DOWN:
            this.threeDEngine.cameraManager.incrementCameraRotate(0, -incrementRotation.y, undefined);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_LEFT:
            this.threeDEngine.cameraManager.incrementCameraRotate(-incrementRotation.x, 0, undefined);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_RIGHT:
            this.threeDEngine.cameraManager.incrementCameraRotate(incrementRotation.x, 0, undefined);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_Z:
            this.threeDEngine.cameraManager.incrementCameraRotate(0, 0, incrementRotation.z);
            break;

          case _CameraControls.cameraControlsActions.ROTATE_MZ:
            this.threeDEngine.cameraManager.incrementCameraRotate(0, 0, -incrementRotation.z);
            break;

          case _CameraControls.cameraControlsActions.ROTATE:
            this.threeDEngine.cameraManager.autoRotate(movieFilter); // movie filter

            break;

          case _CameraControls.cameraControlsActions.ZOOM_IN:
            this.threeDEngine.cameraManager.incrementCameraZoom(-incrementZoom);
            break;

          case _CameraControls.cameraControlsActions.ZOOM_OUT:
            this.threeDEngine.cameraManager.incrementCameraZoom(incrementZoom);
            break;

          case _CameraControls.cameraControlsActions.PAN_HOME:
            this.threeDEngine.cameraManager.resetCamera();
            break;

          case _CameraControls.cameraControlsActions.WIREFRAME:
            this.threeDEngine.setWireframe(!this.threeDEngine.getWireframe());
            break;
        }

        this.threeDEngine.updateControls();
      }
    }
  }, {
    key: "getCanvasElement",
    value: function getCanvasElement() {
      return this.sceneRef && this.sceneRef.current.getElementsByTagName('canvas')[0];
    }
  }, {
    key: "shouldEngineTraverse",
    value: function shouldEngineTraverse() {
      // TODO: check if new instance added, check if split meshes changed?
      return true;
    }
  }, {
    key: "frameResizing",
    value: function frameResizing(width, height, targetRef) {
      this.threeDEngine.resize();
    }
  }, {
    key: "render",
    value: function render() {
      var _this$props3 = this.props,
          classes = _this$props3.classes,
          cameraOptions = _this$props3.cameraOptions,
          captureOptions = _this$props3.captureOptions;
      var cameraControls = cameraOptions.cameraControls;
      var cameraControlsHandler = cameraControls.cameraControlsHandler ? cameraControls.cameraControlsHandler : this.defaultCameraControlsHandler;
      var captureInstance = null;

      if (captureOptions) {
        var captureControls = captureOptions.captureControls;
        var captureControlsHandler = captureControls && captureControls.captureControlsHandler ? captureControls.captureControlsHandler : this.defaultCaptureControlsHandler;
        captureInstance = captureControls && captureControls.instance ? /*#__PURE__*/_react["default"].createElement(captureControls.instance, _extends({
          ref: this.captureControls,
          captureControlsHandler: captureControlsHandler
        }, captureControls.props)) : null;
      }

      return /*#__PURE__*/_react["default"].createElement(_reactResizeDetector["default"], {
        skipOnMount: "true",
        onResize: this.frameResizing
      }, /*#__PURE__*/_react["default"].createElement("div", {
        className: classes.container,
        ref: this.sceneRef
      }, /*#__PURE__*/_react["default"].createElement(cameraControls.instance, _extends({
        ref: this.cameraControls,
        cameraControlsHandler: cameraControlsHandler
      }, cameraControls.props)), captureInstance));
    }
  }]);

  return Canvas;
}(_react.Component);

Canvas.defaultProps = {
  cameraOptions: {
    angle: 50,
    near: 0.01,
    far: 1000,
    baseZoom: 1,
    reset: false,
    autorotate: false,
    wireframe: false,
    depthWrite: true,
    zoomTo: undefined,
    cameraControls: {
      instance: null,
      props: {}
    },
    rotateSpeed: 0.5
  },
  captureOptions: undefined,
  backgroundColor: 0x000000,
  pickingEnabled: true,
  linesThreshold: 2000,
  hoverListeners: [],
  emptyHoverListeners: [],
  threeDObjects: [],
  cameraHandler: function cameraHandler() {},
  selectionStrategy: _SelectionManager.selectionStrategies.nearest,
  onSelection: function onSelection() {},
  setColorHandler: function setColorHandler() {
    return true;
  },
  onMount: function onMount() {},
  modelVersion: 0,
  updateStarted: function updateStarted() {},
  updateEnded: function updateEnded() {},
  dracoDecoderPath: 'https://www.gstatic.com/draco/versioned/decoders/1.5.5/'
};
Canvas.propTypes = {
  /**
   * (Proxy) Instances to visualize
   */
  data: _propTypes["default"].array.isRequired,

  /**
   * Model identifier needed to propagate updates on async changes
   */
  modelVersion: _propTypes["default"].number,

  /**
   * Options to customize camera
   */
  cameraOptions: _propTypes["default"].object,

  /**
   * Options to customize capture features
   */
  captureOptions: _propTypes["default"].shape({
    /**
     * Capture controls component definition
     */
    captureControls: _propTypes["default"].shape({
      /**
       * Component instance
       */
      instance: _propTypes["default"].any,

      /**
       * Component props
       */
      props: _propTypes["default"].shape({})
    }),

    /**
     * Recorder Options
     */
    recorderOptions: _propTypes["default"].shape({
      /**
       * Media Recorder options
       */
      mediaRecorderOptions: _propTypes["default"].shape({
        mimeType: _propTypes["default"].string
      }),
      blobOptions: _propTypes["default"].shape({
        type: _propTypes["default"].string
      })
    }),

    /**
     * Screenshot Options
     */
    screenshotOptions: _propTypes["default"].shape({
      /**
       * A function taking DOM node as argument. Should return true if passed node should be included in the output. Excluding node means excluding it's children as well.
       */
      filter: _propTypes["default"].func,

      /**
       * The pixel ratio of the captured image. Default use the actual pixel ratio of the device. Set 1 to use as initial-scale 1 for the image.
       */
      pixelRatio: _propTypes["default"].number,

      /**
       * A number between 0 and 1 indicating image quality (e.g. 0.92 => 92%) of the JPEG image.
       */
      quality: _propTypes["default"].number,

      /**
       * Screenshot desired resolution
       */
      resolution: _propTypes["default"].shape({
        height: _propTypes["default"].number.isRequired,
        width: _propTypes["default"].number.isRequired
      })
    }).isRequired
  }),

  /**
   * Three JS objects to add to the scene
   */
  threeDObjects: _propTypes["default"].array,

  /**
   * Function to callback on camera changes
   */
  cameraHandler: _propTypes["default"].func,

  /**
   * function to apply the selection strategy
   */
  selectionStrategy: _propTypes["default"].func,

  /**
   * Function to callback on selection changes
   */
  onSelection: _propTypes["default"].func,

  /**
   * Function to callback on set color changes. Return true to apply default behavior after or false otherwise
   */
  setColorHandler: _propTypes["default"].func,

  /**
   * Function to callback on component did mount with scene
   */
  onMount: _propTypes["default"].func,

  /**
   * Scene background color
   */
  backgroundColor: _propTypes["default"].number,

  /**
   * Boolean to enable/disable 3d picking
   */
  pickingEnabled: _propTypes["default"].bool,

  /**
   * Threshold to limit scene complexity
   */
  linesThreshold: _propTypes["default"].number,

  /**
   * Array of hover handlers to callback
   */
  hoverListeners: _propTypes["default"].array,

  /**
   * Array of hover handlers to callback
   */
  emptyHoverListener: _propTypes["default"].func,

  /**
   * Function to callback when the loading of elements of the canvas starts
   */
  updateStarted: _propTypes["default"].func,

  /**
   * Function to callback when the loading of elements of the canvas ends
   */
  updateEnded: _propTypes["default"].func,

  /**
   * Path to the draco decoder
   */
  dracoDecoderPath: _propTypes["default"].string
};

var _default = (0, _core.withStyles)(styles)(Canvas);

exports["default"] = _default;