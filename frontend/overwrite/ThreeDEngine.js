"use strict";

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;

var THREE = _interopRequireWildcard(require("three"));

var _RenderPass = require("three/examples/jsm/postprocessing/RenderPass.js");

var _EffectComposer = require("three/examples/jsm/postprocessing/EffectComposer.js");

var _ShaderPass = require("three/examples/jsm/postprocessing/ShaderPass.js");

var _BloomPass = require("three/examples/jsm/postprocessing/BloomPass.js");

var _FilmPass = require("three/examples/jsm/postprocessing/FilmPass.js");

var _FocusShader = require("three/examples/jsm/shaders/FocusShader.js");

var _MeshFactory = _interopRequireDefault(require("./MeshFactory"));

var _CameraManager = _interopRequireDefault(require("./CameraManager"));

var _Instance = _interopRequireDefault(require("@metacell/geppetto-meta-core/model/Instance"));

var _ArrayInstance = _interopRequireDefault(require("@metacell/geppetto-meta-core//model/ArrayInstance"));

var _Type = _interopRequireDefault(require("@metacell/geppetto-meta-core/model/Type"));

var _Variable = _interopRequireDefault(require("@metacell/geppetto-meta-core/model/Variable"));

var _SimpleInstance = _interopRequireDefault(require("@metacell/geppetto-meta-core/model/SimpleInstance"));

var _util = require("./util");

var _Utility = require("@metacell/geppetto-meta-core/Utility");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _getRequireWildcardCache(nodeInterop) { if (typeof WeakMap !== "function") return null; var cacheBabelInterop = new WeakMap(); var cacheNodeInterop = new WeakMap(); return (_getRequireWildcardCache = function _getRequireWildcardCache(nodeInterop) { return nodeInterop ? cacheNodeInterop : cacheBabelInterop; })(nodeInterop); }

function _interopRequireWildcard(obj, nodeInterop) { if (!nodeInterop && obj && obj.__esModule) { return obj; } if (obj === null || _typeof(obj) !== "object" && typeof obj !== "function") { return { "default": obj }; } var cache = _getRequireWildcardCache(nodeInterop); if (cache && cache.has(obj)) { return cache.get(obj); } var newObj = {}; var hasPropertyDescriptor = Object.defineProperty && Object.getOwnPropertyDescriptor; for (var key in obj) { if (key !== "default" && Object.prototype.hasOwnProperty.call(obj, key)) { var desc = hasPropertyDescriptor ? Object.getOwnPropertyDescriptor(obj, key) : null; if (desc && (desc.get || desc.set)) { Object.defineProperty(newObj, key, desc); } else { newObj[key] = obj[key]; } } } newObj["default"] = obj; if (cache) { cache.set(obj, newObj); } return newObj; }

function _createForOfIteratorHelper(o, allowArrayLike) { var it = typeof Symbol !== "undefined" && o[Symbol.iterator] || o["@@iterator"]; if (!it) { if (Array.isArray(o) || (it = _unsupportedIterableToArray(o)) || allowArrayLike && o && typeof o.length === "number") { if (it) o = it; var i = 0; var F = function F() {}; return { s: F, n: function n() { if (i >= o.length) return { done: true }; return { done: false, value: o[i++] }; }, e: function e(_e) { throw _e; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var normalCompletion = true, didErr = false, err; return { s: function s() { it = it.call(o); }, n: function n() { var step = it.next(); normalCompletion = step.done; return step; }, e: function e(_e2) { didErr = true; err = _e2; }, f: function f() { try { if (!normalCompletion && it["return"] != null) it["return"](); } finally { if (didErr) throw err; } } }; }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function asyncGeneratorStep(gen, resolve, reject, _next, _throw, key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { Promise.resolve(value).then(_next, _throw); } }

function _asyncToGenerator(fn) { return function () { var self = this, args = arguments; return new Promise(function (resolve, reject) { var gen = fn.apply(self, args); function _next(value) { asyncGeneratorStep(gen, resolve, reject, _next, _throw, "next", value); } function _throw(err) { asyncGeneratorStep(gen, resolve, reject, _next, _throw, "throw", err); } _next(undefined); }); }; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) { symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); } keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

require('./TrackballControls');

var ThreeDEngine = /*#__PURE__*/function () {
  function ThreeDEngine(containerRef, cameraOptions, cameraHandler, captureOptions, onSelection, backgroundColor, pickingEnabled, linesThreshold, hoverListeners, emptyHoverListener, setColorHandler, selectionStrategy, updateStarted, updateEnded, dracoDecoderPath) {
    _classCallCheck(this, ThreeDEngine);

    this.scene = new THREE.Scene();
    this.scene.background = new THREE.Color(backgroundColor);
    this.cameraManager = null;
    this.renderer = null;
    this.controls = null;
    this.mouse = {
      x: 0,
      y: 0
    };
    this.mouseContainer = {
      x: 0,
      y: 0
    };
    this.frameId = null;
    this.meshFactory = new _MeshFactory["default"](this.scene, linesThreshold, cameraOptions.depthWrite, 300, 1, null, dracoDecoderPath, null);
    this.pickingEnabled = pickingEnabled;
    this.hoverListeners = hoverListeners;
    this.emptyHoverListener = emptyHoverListener;
    this.cameraHandler = cameraHandler;
    this.setColorHandler = setColorHandler;
    this.selectionStrategy = selectionStrategy;
    this.containerRef = containerRef;
    this.width = containerRef.clientWidth;
    this.height = containerRef.clientHeight;
    this.lastRequestFrame = 0;
    this.lastRenderTimer = new Date();
    this.updateStarted = updateStarted;
    this.updateEnded = updateEnded; // Setup Camera

    this.setupCamera(cameraOptions, this.width / this.height); // Setup Renderer

    this.setupRenderer(containerRef, {
      antialias: true,
      alpha: true,
      preserveDrawingBuffer: captureOptions !== undefined
    }); // Setup Lights

    this.setupLights(); // Setup Controls

    this.setupControls(); // Setup Listeners

    this.setupListeners(onSelection);
    this.instancesMap = new Map();
    this.start = this.start.bind(this);
    this.stop = this.stop.bind(this);
    this.animate = this.animate.bind(this);
    this.renderScene = this.renderScene.bind(this);
    this.resize = this.resize.bind(this);
  }
  /**
   * Setups the camera
   * @param cameraOptions
   * @param aspect
   */


  _createClass(ThreeDEngine, [{
    key: "setupCamera",
    value: function setupCamera(cameraOptions, aspect) {
      this.cameraManager = new _CameraManager["default"](this, _objectSpread(_objectSpread({}, cameraOptions), {}, {
        aspect: aspect
      }));
    }
    /**
     * Setups the renderer
     * @param containerRef
     */

  }, {
    key: "setupRenderer",
    value: function setupRenderer(containerRef, options) {
      this.renderer = new THREE.WebGLRenderer(options);
      this.renderer.setSize(this.width, this.height);
      this.renderer.setPixelRatio(window.devicePixelRatio);
      this.renderer.autoClear = false;
      containerRef.appendChild(this.renderer.domElement);
      this.configureRenderer(false);
    }
    /**
     *
     * @param shaders
     */

  }, {
    key: "configureRenderer",
    value: function configureRenderer(shaders) {
      if (shaders === undefined) {
        shaders = false;
      }

      var renderModel = new _RenderPass.RenderPass(this.scene, this.cameraManager.getCamera());
      this.composer = new _EffectComposer.EffectComposer(this.renderer);

      if (shaders) {
        var effectBloom = new _BloomPass.BloomPass(0.75); // todo: grayscale shouldn't be false

        var effectFilm = new _FilmPass.FilmPass(0.5, 0.5, 1448, false);
        var effectFocus = new _ShaderPass.ShaderPass(_FocusShader.FocusShader);
        effectFocus.uniforms['screenWidth'].value = this.width;
        effectFocus.uniforms['screenHeight'].value = this.height;
        effectFocus.renderToScreen = true;
        this.composer.addPass(renderModel);
        this.composer.addPass(effectBloom);
        this.composer.addPass(effectFilm);
        this.composer.addPass(effectFocus);
      } else {
        // standard
        renderModel.renderToScreen = true;
        this.composer.addPass(renderModel);
      }
    }
    /**
     * Setups the lights
     */

  }, {
    key: "setupLights",
    value: function setupLights() {
      var ambientLight = new THREE.AmbientLight(0x0c0c0c);
      this.scene.add(ambientLight);
      var spotLight = new THREE.SpotLight(0xffffff);
      spotLight.position.set(-30, 60, 60);
      spotLight.castShadow = true;
      this.scene.add(spotLight);
      this.cameraManager.getCamera().add(new THREE.PointLight(0xffffff, 1));
    }
  }, {
    key: "setupControls",
    value: function setupControls() {
      this.controls = new THREE.TrackballControls(this.cameraManager.getCamera(), this.renderer.domElement, this.cameraHandler);
      this.controls.noZoom = false;
      this.controls.noPan = false;
    }
    /**
     * Returns intersected objects from mouse click
     *
     * @returns {Array} a list of objects intersected by the current mouse coordinates
     */

  }, {
    key: "getIntersectedObjects",
    value: function getIntersectedObjects() {
      // create a Ray with origin at the mouse position and direction into th scene (camera direction)
      var vector = new THREE.Vector3(this.mouse.x, this.mouse.y, 1);
      vector.unproject(this.cameraManager.getCamera());
      var raycaster = new THREE.Raycaster(this.cameraManager.getCamera().position, vector.sub(this.cameraManager.getCamera().position).normalize());
      raycaster.linePrecision = this.meshFactory.getLinePrecision();
      var visibleChildren = [];
      this.scene.traverse(function (child) {
        if (child.visible && !(child.clickThrough === true)) {
          if (child.geometry != null) {
            if (child.type !== 'Points') {
              child.geometry.computeBoundingBox();
            }

            visibleChildren.push(child);
          }
        }
      }); // returns an array containing all objects in the scene with which the ray intersects

      return raycaster.intersectObjects(visibleChildren);
    }
    /**
     * Adds instances to the ThreeJS Scene
     * @param proxyInstances
     */

  }, {
    key: "addInstancesToScene",
    value: function () {
      var _addInstancesToScene = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee(proxyInstances) {
        var _this2, i, pInstance, geppettoInstance;
        return regeneratorRuntime.wrap(function _callee$(_context) {
          while (1) {
            switch (_context.prev = _context.next) {
              case 0:
                _this2 = this;
                this.instancesMap.clear();
                if (Array.isArray(proxyInstances)) {
                  for (i = 0; i < proxyInstances.length; i++) {
                    pInstance = proxyInstances[i];
                    geppettoInstance = Instances.getInstance(pInstance.instancePath);
                    if (geppettoInstance) {
                      _this2.traverseAndMapInstance(pInstance, geppettoInstance);
                    }
                  }
                }
                _context.next = 5;
                return this.meshFactory.start(this.instancesMap);

              case 5:
                this.updateGroupMeshes(proxyInstances);

              case 6:
              case "end":
                return _context.stop();
            }
          }
        }, _callee, this);
      }));

      function addInstancesToScene(_x) {
        return _addInstancesToScene.apply(this, arguments);
      }

      return addInstancesToScene;
    }()
  }, {
    key: "traverseAndMapInstance",
    value: function traverseAndMapInstance(proxyInstance, geppettoInstance) {
      try {
        if ((0, _util.hasVisualValue)(geppettoInstance)) {
          this.instancesMap.set(geppettoInstance.getInstancePath(), proxyInstance);
        } else if ((0, _util.hasVisualType)(geppettoInstance)) {
          if (geppettoInstance.getType().getMetaType() !== GEPPETTO.Resources.ARRAY_TYPE_NODE && geppettoInstance.getVisualType()) {
            this.instancesMap.set(geppettoInstance.getInstancePath(), proxyInstance);
          }
          if (geppettoInstance.getMetaType() === GEPPETTO.Resources.INSTANCE_NODE) {
            var children = geppettoInstance.getChildren();
            for (var i = 0; i < children.length; i++) {
              this.traverseAndMapInstance(proxyInstance, children[i]);
            }
          } else if (geppettoInstance.getMetaType() === GEPPETTO.Resources.ARRAY_INSTANCE_NODE) {
            for (var j = 0; j < geppettoInstance.length; j++) {
              this.traverseAndMapInstance(proxyInstance, geppettoInstance[j]);
            }
          }
        }
      } catch (e) {
        console.error(e);
      }
    }
    /*
     * Check that the material for the already present instance did not change.
     * return true if the color changed, otherwise false.
     */

  }, {
    key: "checkMaterial",
    value: function checkMaterial(mesh, instance) {
      if (mesh.type === 'Mesh') {
        var _instance$color, _instance$color2, _instance$color3, _instance$color4;

        if (mesh.material.color.r === (instance === null || instance === void 0 ? void 0 : (_instance$color = instance.color) === null || _instance$color === void 0 ? void 0 : _instance$color.r) && mesh.material.color.g === (instance === null || instance === void 0 ? void 0 : (_instance$color2 = instance.color) === null || _instance$color2 === void 0 ? void 0 : _instance$color2.g) && mesh.material.color.b === (instance === null || instance === void 0 ? void 0 : (_instance$color3 = instance.color) === null || _instance$color3 === void 0 ? void 0 : _instance$color3.b) && mesh.material.color.opacity === (instance === null || instance === void 0 ? void 0 : (_instance$color4 = instance.color) === null || _instance$color4 === void 0 ? void 0 : _instance$color4.a)) {
          return false;
        } else {
          return true;
        }
      } else if (mesh.type === 'Group') {
        var changed = false;

        var _iterator = _createForOfIteratorHelper(mesh.children),
            _step;

        try {
          for (_iterator.s(); !(_step = _iterator.n()).done;) {
            var child = _step.value;

            if (this.checkMaterial(child, instance)) {
              changed = true;
            }
          }
        } catch (err) {
          _iterator.e(err);
        } finally {
          _iterator.f();
        }

        return changed;
      }
    }
  }, {
    key: "updateInstanceMaterial",
    value: function updateInstanceMaterial(mesh, instance) {
      var _iterator2 = _createForOfIteratorHelper(this.scene.children),
          _step2;

      try {
        for (_iterator2.s(); !(_step2 = _iterator2.n()).done;) {
          var child = _step2.value;

          if (child.instancePath === mesh.instancePath && child.uuid === mesh.uuid) {
            if ((instance === null || instance === void 0 ? void 0 : instance.color) !== undefined) {
              this.setInstanceMaterial(child, instance);
              break;
            } else {
              instance.color = GEPPETTO.Resources.COLORS.DEFAULT;
              this.setInstanceMaterial(child, instance);
            }
          }
        }
      } catch (err) {
        _iterator2.e(err);
      } finally {
        _iterator2.f();
      }
    }
  }, {
    key: "setInstanceMaterial",
    value: function setInstanceMaterial(mesh, instance) {
      if (mesh.type === 'Mesh') {
        this.meshFactory.setThreeColor(mesh.material.color, instance.color);

        if (instance.color.a) {
          mesh.material.transparent = true;
          mesh.material.opacity = instance.color.a;
        }
      } else if (mesh.type === 'Group') {
        var _iterator3 = _createForOfIteratorHelper(mesh.children),
            _step3;

        try {
          for (_iterator3.s(); !(_step3 = _iterator3.n()).done;) {
            var child = _step3.value;
            this.setInstanceMaterial(child, instance);
          }
        } catch (err) {
          _iterator3.e(err);
        } finally {
          _iterator3.f();
        }
      }
    }
    /**
     * Clears the scene
     *
     * we have the list of strings instances
     * we have the global Instances from the model
     * we have the obj instances in the threeJS scene
     */

  }, {
    key: "checkInstanceToRemove",
    value: function checkInstanceToRemove(geppettoInstance, proxyInstance, toRemove, pathsToRemove) {
      try {
        if ((0, _util.hasVisualValue)(geppettoInstance)) {
          var geppettoIndex = pathsToRemove.indexOf(geppettoInstance.getInstancePath());

          if (geppettoIndex > -1) {
            if (this.checkMaterial(toRemove[geppettoIndex], proxyInstance)) {
              this.updateInstanceMaterial(toRemove[geppettoIndex], proxyInstance);
            }

            toRemove.splice(geppettoIndex, 1);
            pathsToRemove.splice(geppettoIndex, 1);
            return true;
          }

          return false;
        } else if ((0, _util.hasVisualType)(geppettoInstance)) {
          if (geppettoInstance.getType().getMetaType() !== GEPPETTO.Resources.ARRAY_TYPE_NODE && geppettoInstance.getVisualType()) {
            var geppettoIndex = pathsToRemove.indexOf(geppettoInstance.getInstancePath());

            if (geppettoIndex > -1) {
              if (this.checkMaterial(toRemove[geppettoIndex], proxyInstance)) {
                this.updateInstanceMaterial(toRemove[geppettoIndex], proxyInstance);
              }

              toRemove.splice(geppettoIndex, 1);
              pathsToRemove.splice(geppettoIndex, 1);
              return true;
            }

            return false;
          } // this block keeps traversing the instances


          if (geppettoInstance.getMetaType() === GEPPETTO.Resources.INSTANCE_NODE) {
            var returnValue = false;
            var children = geppettoInstance.getChildren();

            for (var i = 0; i < children.length; i++) {
              var instanceReturn = this.checkInstanceToRemove(children[i], proxyInstance, toRemove, pathsToRemove);
              returnValue = returnValue || instanceReturn;
            }

            return returnValue;
          } else if (geppettoInstance.getMetaType() === GEPPETTO.Resources.ARRAY_INSTANCE_NODE) {
            var returnValue = false;

            for (var _i = 0; _i < geppettoInstance.length; _i++) {
              var _instanceReturn = this.checkInstanceToRemove(geppettoInstance[_i], proxyInstance, toRemove, pathsToRemove);

              returnValue = returnValue || _instanceReturn;
            }

            return returnValue;
          }
        }
      } catch (e) {
        console.error(e);
      }
    }
  }, {
    key: "clearScene",
    value: function () {
      var _clearScene = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee2(proxyInstances) {
        var pathsToRemove, sortedInstances, toRemove, i, _sortedInstances$i, geppettoInstance, _check, _iterator4, _step4, child;

        return regeneratorRuntime.wrap(function _callee2$(_context2) {
          while (1) {
            switch (_context2.prev = _context2.next) {
              case 0:
                pathsToRemove = [];
                sortedInstances = [];
                toRemove = this.scene.children.filter(function (child) {
                  if (child.type === 'Mesh' || child.type === 'Group') {
                    pathsToRemove.push(child.instancePath);
                    return true;
                  }

                  return false;
                });

                if (!proxyInstances) {
                  _context2.next = 23;
                  break;
                }

                sortedInstances = proxyInstances.sort(function (a, b) {
                  if (a.instancePath < b.instancePath) {
                    return -1;
                  }

                  if (a.instancePath > b.instancePath) {
                    return 1;
                  }

                  return 0;
                });

                if (!(toRemove.length === 0)) {
                  _context2.next = 7;
                  break;
                }

                return _context2.abrupt("return", sortedInstances);

              case 7:
                i = sortedInstances.length - 1;

              case 8:
                if (!(i >= 0)) {
                  _context2.next = 21;
                  break;
                }

                geppettoInstance = Instances.getInstance((_sortedInstances$i = sortedInstances[i]) === null || _sortedInstances$i === void 0 ? void 0 : _sortedInstances$i.instancePath);

                if (!geppettoInstance) {
                  _context2.next = 17;
                  break;
                }

                _context2.next = 13;
                return this.checkInstanceToRemove(geppettoInstance, sortedInstances[i], toRemove, pathsToRemove);

              case 13:
                _check = _context2.sent;

                if (_check) {
                  sortedInstances.splice(i, 1);
                }

                _context2.next = 18;
                break;

              case 17:
                sortedInstances.splice(i, 1);

              case 18:
                i--;
                _context2.next = 8;
                break;

              case 21:
                _context2.next = 25;
                break;

              case 23:
                console.error("Give me an empty list if you want to wipe all the instances from the 3d viewer.");
                return _context2.abrupt("return", []);

              case 25:
                _iterator4 = _createForOfIteratorHelper(toRemove);

                try {
                  for (_iterator4.s(); !(_step4 = _iterator4.n()).done;) {
                    child = _step4.value;
                    this.meshFactory.cleanWith3DObject(child);
                    this.scene.remove(child);
                  }
                } catch (err) {
                  _iterator4.e(err);
                } finally {
                  _iterator4.f();
                }

                return _context2.abrupt("return", sortedInstances);

              case 28:
              case "end":
                return _context2.stop();
            }
          }
        }, _callee2, this);
      }));

      function clearScene(_x2) {
        return _clearScene.apply(this, arguments);
      }

      return clearScene;
    }()
  }, {
    key: "updateInstancesColor",
    value: function updateInstancesColor(proxyInstances) {
      var sortedInstances = proxyInstances.sort(function (a, b) {
        if (a.instancePath < b.instancePath) {
          return -1;
        }

        if (a.instancePath > b.instancePath) {
          return 1;
        }

        return 0;
      });

      var _iterator5 = _createForOfIteratorHelper(sortedInstances),
          _step5;

      try {
        for (_iterator5.s(); !(_step5 = _iterator5.n()).done;) {
          var pInstance = _step5.value;

          if (pInstance.color) {
            this.setInstanceColor(pInstance.instancePath, pInstance.color);
          }

          if (pInstance.visualGroups) {
            var instance = Instances.getInstance(pInstance.instancePath);
            var visualGroups = this.getVisualElements(instance, pInstance.visualGroups);
            this.setSplitGroupsColor(pInstance.instancePath, visualGroups);
          }
        }
      } catch (err) {
        _iterator5.e(err);
      } finally {
        _iterator5.f();
      }
    }
  }, {
    key: "updateInstancesConnectionLines",
    value: function updateInstancesConnectionLines(proxyInstances) {
      var _iterator6 = _createForOfIteratorHelper(proxyInstances),
          _step6;

      try {
        for (_iterator6.s(); !(_step6 = _iterator6.n()).done;) {
          var pInstance = _step6.value;
          var mode = pInstance.showConnectionLines ? pInstance.showConnectionLines : false;
          this.showConnectionLines(pInstance.instancePath, mode);
        }
      } catch (err) {
        _iterator6.e(err);
      } finally {
        _iterator6.f();
      }
    }
    /**
     * Sets the color of the instances
     *
     * @param path
     * @param color
     */

  }, {
    key: "setInstanceColor",
    value: function setInstanceColor(path, color) {
      var entity = Instances.getInstance(path);

      if (entity && this.setColorHandler(entity)) {
        if (entity instanceof _Instance["default"] || entity instanceof _ArrayInstance["default"] || entity instanceof _SimpleInstance["default"]) {
          this.meshFactory.setColor(path, color);

          if (typeof entity.getChildren === 'function') {
            var children = entity.getChildren();

            for (var i = 0; i < children.length; i++) {
              this.setInstanceColor(children[i].getInstancePath(), color);
            }
          }
        } else if (entity instanceof _Type["default"] || entity instanceof _Variable["default"]) {
          // fetch all instances for the given type or variable and call hide on each
          var instances = GEPPETTO.ModelFactory.getAllInstancesOf(entity);

          for (var j = 0; j < instances.length; j++) {
            this.setInstanceColor(instances[j].getInstancePath(), color);
          }
        }
      }
    }
    /**
     *
     * @param instancePath
     * @param visualGroups
     */

  }, {
    key: "setSplitGroupsColor",
    value: function setSplitGroupsColor(instancePath, visualGroups) {
      for (var g in visualGroups) {
        // retrieve visual group object
        var group = visualGroups[g]; // get full group name to access group mesh

        var groupName = g;

        if (groupName.indexOf(instancePath) <= -1) {
          groupName = instancePath + '.' + g;
        } // get group mesh


        var groupMesh = this.meshFactory.getMeshes()[groupName];
        groupMesh.visible = true;
        this.meshFactory.setThreeColor(groupMesh.material.color, group.color);
      }
    }
  }, {
    key: "updateGroupMeshes",
    value: function updateGroupMeshes(proxyInstances) {
      var _iterator7 = _createForOfIteratorHelper(proxyInstances),
          _step7;

      try {
        for (_iterator7.s(); !(_step7 = _iterator7.n()).done;) {
          var pInstance = _step7.value;

          if (pInstance.visualGroups) {
            var instance = Instances.getInstance(pInstance.instancePath);
            var visualGroups = this.getVisualElements(instance, pInstance.visualGroups);
            this.meshFactory.splitGroups(instance, visualGroups);
          }
        }
      } catch (err) {
        _iterator7.e(err);
      } finally {
        _iterator7.f();
      }

      var meshes = this.meshFactory.getMeshes();

      for (var meshKey in meshes) {
        this.addToScene(meshes[meshKey]);
      }
    }
  }, {
    key: "getVisualElements",
    value: function getVisualElements(instance, visualGroups) {
      var groups = {};

      if (visualGroups.index != null) {
        var vg = instance.getVisualGroups()[visualGroups.index];
        var visualElements = vg.getVisualGroupElements();
        var allElements = [];

        for (var i = 0; i < visualElements.length; i++) {
          if (visualElements[i].getValue() != null) {
            allElements.push(visualElements[i].getValue());
          }
        }

        var minDensity = Math.min.apply(null, allElements);
        var maxDensity = Math.max.apply(null, allElements); // highlight all reference nodes

        for (var j = 0; j < visualElements.length; j++) {
          groups[visualElements[j].getId()] = {};
          var color = visualElements[j].getColor();

          if (visualElements[j].getValue() != null) {
            var intensity = 1;

            if (maxDensity !== minDensity) {
              intensity = (visualElements[j].getValue() - minDensity) / (maxDensity - minDensity);
            }

            color = (0, _Utility.rgbToHex)(255, Math.floor(255 - 255 * intensity), 0);
          }

          groups[visualElements[j].getId()].color = color;
        }
      }

      for (var c in visualGroups.custom) {
        if (c in groups) {
          groups[c].color = visualGroups.custom[c].color;
        }
      }

      return groups;
    }
    /**
     * Show connection lines for this instance.
     *
     * @param instancePath
     * @param {boolean} mode - Show or hide connection lines
     */

  }, {
    key: "showConnectionLines",
    value: function showConnectionLines(instancePath, mode) {
      if (mode == null) {
        mode = true;
      }

      var entity = Instances.getInstance(instancePath);

      if (entity instanceof _Instance["default"] || entity instanceof _ArrayInstance["default"]) {
        // show or hide connection lines
        if (mode) {
          this.showConnectionLinesForInstance(entity);
        } else {
          this.removeConnectionLines(entity);
        }
      } else if (entity instanceof _Type["default"] || entity instanceof _Variable["default"]) {
        // fetch all instances for the given type or variable and call hide on each
        var instances = GEPPETTO.ModelFactory.getAllInstancesOf(entity);

        for (var j = 0; j < instances.length; j++) {
          if ((0, _util.hasVisualType)(instances[j])) {
            this.showConnectionLines(instances[j], mode);
          }
        }
      }
    }
    /**
     *
     *
     * @param instance
     */

  }, {
    key: "showConnectionLinesForInstance",
    value: function showConnectionLinesForInstance(instance) {
      var connections = instance.getConnections();
      var mesh = this.meshFactory.meshes[instance.getInstancePath()];
      var inputs = {};
      var outputs = {};
      var defaultOrigin = mesh.position.clone();

      for (var c = 0; c < connections.length; c++) {
        var connection = connections[c];
        var type = connection.getA().getPath() === instance.getInstancePath() ? GEPPETTO.Resources.OUTPUT : GEPPETTO.Resources.INPUT;
        var thisEnd = connection.getA().getPath() === instance.getInstancePath() ? connection.getA() : connection.getB();
        var otherEnd = connection.getA().getPath() === instance.getInstancePath() ? connection.getB() : connection.getA();
        var otherEndPath = otherEnd.getPath();
        var otherEndMesh = this.meshFactory.meshes[otherEndPath];
        var destination = void 0;
        var origin = void 0;

        if (thisEnd.getPoint() === undefined) {
          // same as before
          origin = defaultOrigin;
        } else {
          // the specified coordinate
          var p = thisEnd.getPoint();
          origin = new THREE.Vector3(p.x + mesh.position.x, p.y + mesh.position.y, p.z + mesh.position.z);
        }

        if (otherEnd.getPoint() === undefined) {
          // same as before
          destination = otherEndMesh.position.clone();
        } else {
          // the specified coordinate
          var _p = otherEnd.getPoint();

          destination = new THREE.Vector3(_p.x + otherEndMesh.position.x, _p.y + otherEndMesh.position.y, _p.z + otherEndMesh.position.z);
        }

        var geometry = new THREE.Geometry();
        geometry.vertices.push(origin, destination);
        geometry.verticesNeedUpdate = true;
        geometry.dynamic = true;
        var colour = null;

        if (type === GEPPETTO.Resources.INPUT) {
          colour = GEPPETTO.Resources.COLORS.INPUT_TO_SELECTED; // figure out if connection is both, input and output

          if (outputs[otherEndPath]) {
            colour = GEPPETTO.Resources.COLORS.INPUT_AND_OUTPUT;
          }

          if (inputs[otherEndPath]) {
            inputs[otherEndPath].push(connection.getInstancePath());
          } else {
            inputs[otherEndPath] = [];
            inputs[otherEndPath].push(connection.getInstancePath());
          }
        } else if (type === GEPPETTO.Resources.OUTPUT) {
          colour = GEPPETTO.Resources.COLORS.OUTPUT_TO_SELECTED; // figure out if connection is both, input and output

          if (inputs[otherEndPath]) {
            colour = GEPPETTO.Resources.COLORS.INPUT_AND_OUTPUT;
          }

          if (outputs[otherEndPath]) {
            outputs[otherEndPath].push(connection.getInstancePath());
          } else {
            outputs[otherEndPath] = [];
            outputs[otherEndPath].push(connection.getInstancePath());
          }
        }

        var material = new THREE.LineDashedMaterial({
          dashSize: 3,
          gapSize: 1
        });
        this.meshFactory.setThreeColor(material.color, colour);
        var line = new THREE.LineSegments(geometry, material);
        line.updateMatrixWorld(true);

        if (this.meshFactory.connectionLines[connection.getInstancePath()]) {
          this.scene.remove(this.meshFactory.connectionLines[connection.getInstancePath()]);
        }

        this.scene.add(line);
        this.meshFactory.connectionLines[connection.getInstancePath()] = line;
      }
    }
    /**
     * Removes connection lines, all if nothing is passed in or just the ones passed in.
     *
     * @param instance - optional, instance for which we want to remove the connections
     */

  }, {
    key: "removeConnectionLines",
    value: function removeConnectionLines(instance) {
      if (instance !== undefined) {
        var connections = instance.getConnections(); // get connections for given instance and remove only those

        var lines = this.meshFactory.connectionLines;

        for (var i = 0; i < connections.length; i++) {
          if (Object.prototype.hasOwnProperty.call(lines, connections[i].getInstancePath())) {
            // remove the connection line from the scene
            this.scene.remove(lines[connections[i].getInstancePath()]); // remove the conneciton line from the GEPPETTO list of connection lines

            delete lines[connections[i].getInstancePath()];
          }
        }
      } else {
        // remove all connection lines
        var _lines = this.meshFactory.connectionLines;

        for (var key in _lines) {
          if (Object.prototype.hasOwnProperty.call(_lines, key)) {
            this.scene.remove(_lines[key]);
          }
        }

        this.meshFactory.connectionLines = [];
      }
    }
    /**
     * Set up the listeners use to detect mouse movement and window resizing
     */

  }, {
    key: "setupListeners",
    value: function setupListeners(onSelection) {
      var that = this;
      this.controls.addEventListener('start', function (e) {
        that.requestFrame();
      });
      this.controls.addEventListener('change', function (e) {
        that.requestFrame();
      });
      this.controls.addEventListener('stop', function (e) {
        that.stop();
      }); // when the mouse moves, call the given function

      this.renderer.domElement.addEventListener('mousedown', function (event) {
        that.clientX = event.clientX;
        that.clientY = event.clientY;
      }, false); // when the mouse moves, call the given function

      this.renderer.domElement.addEventListener('mouseup', function (event) {
        if (event.target === that.renderer.domElement) {
          var x = event.clientX;
          var y = event.clientY; // If the mouse moved since the mousedown then don't consider this a selection

          if (typeof that.clientX === 'undefined' || typeof that.clientY === 'undefined' || x !== that.clientX || y !== that.clientY) {
            return;
          }

          that.mouse.y = -((event.clientY - that.renderer.domElement.getBoundingClientRect().top) * window.devicePixelRatio / that.renderer.domElement.height) * 2 + 1;
          that.mouse.x = (event.clientX - that.renderer.domElement.getBoundingClientRect().left) * window.devicePixelRatio / that.renderer.domElement.width * 2 - 1;

          if (that.pickingEnabled) {
            var intersects = that.getIntersectedObjects();

            if (intersects.length > 0) {
              // sort intersects
              var compare = function compare(a, b) {
                if (a.distance < b.distance) {
                  return -1;
                }

                if (a.distance > b.distance) {
                  return 1;
                }

                return 0;
              };

              intersects.sort(compare);
              var selectedMap = {}; // Iterate and get the first visible item (they are now ordered by proximity)

              for (var i = 0; i < intersects.length; i++) {
                // figure out if the entity is visible
                var instancePath = '';
                var geometryIdentifier = '';

                if (Object.prototype.hasOwnProperty.call(intersects[i].object, 'instancePath')) {
                  instancePath = intersects[i].object.instancePath;
                  geometryIdentifier = intersects[i].object.geometryIdentifier;
                } else {
                  // weak assumption: if the object doesn't have an instancePath its parent will
                  instancePath = intersects[i].object.parent.instancePath;
                  geometryIdentifier = intersects[i].object.parent.geometryIdentifier;
                }

                if (instancePath != null && Object.prototype.hasOwnProperty.call(that.meshFactory.meshes, instancePath) || Object.prototype.hasOwnProperty.call(that.meshFactory.splitMeshes, instancePath)) {
                  if (geometryIdentifier === undefined) {
                    geometryIdentifier = '';
                  }

                  if (!(instancePath in selectedMap)) {
                    selectedMap[instancePath] = _objectSpread(_objectSpread({}, intersects[i]), {}, {
                      geometryIdentifier: geometryIdentifier,
                      distanceIndex: i
                    });
                  }
                }
              }

              that.requestFrame();
              onSelection(that.selectionStrategy(selectedMap), event);
            }
          }
        }
      }, false);
      this.renderer.domElement.addEventListener('mousemove', function (event) {
        that.mouse.y = -((event.clientY - that.renderer.domElement.getBoundingClientRect().top) * window.devicePixelRatio / that.renderer.domElement.height) * 2 + 1;
        that.mouse.x = (event.clientX - that.renderer.domElement.getBoundingClientRect().left) * window.devicePixelRatio / that.renderer.domElement.width * 2 - 1;
        that.mouseContainer.x = event.clientX;
        that.mouseContainer.y = event.clientY;

        if (that.hoverListeners && that.hoverListeners.length > 0) {
          var intersects = that.getIntersectedObjects();

          if (intersects.length !== 0) {
            for (var listener in that.hoverListeners) {
              that.hoverListeners[listener](intersects, that.mouseContainer.x, that.mouseContainer.y);
            }
          } else {
            that.emptyHoverListener();
          }
        }

      }, false);
    }
    /**
     * Sets whether to use wireframe for the materials of the meshes
     * @param wireframe
     */

  }, {
    key: "setWireframe",
    value: function setWireframe(wireframe) {
      this.wireframe = wireframe;
      var that = this;
      this.scene.traverse(function (child) {
        if (child instanceof THREE.Mesh) {
          if (!(child.material.nowireframe === true)) {
            child.material.wireframe = that.wireframe;
          }
        }
      });
    }
  }, {
    key: "setBackgroundColor",
    value: function setBackgroundColor(color) {
      this.scene.background.getHex();
      var newColor = new THREE.Color(color);

      if (this.scene.background.getHex() !== newColor.getHex()) {
        this.scene.background = newColor;
      }
    }
  }, {
    key: "update",
    value: function () {
      var _update = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee3(proxyInstances, cameraOptions, threeDObjects, toTraverse, newBackgroundColor) {
        var _this = this;

        return regeneratorRuntime.wrap(function _callee3$(_context3) {
          while (1) {
            switch (_context3.prev = _context3.next) {
              case 0:
                this.updateStarted();
                this.setBackgroundColor(newBackgroundColor);
                _context3.next = 4;
                return this.clearScene(proxyInstances);

              case 4:
                proxyInstances = _context3.sent;

                if (!toTraverse) {
                  _context3.next = 12;
                  break;
                }

                _context3.next = 8;
                return this.addInstancesToScene(proxyInstances);

              case 8:
                threeDObjects.forEach(function (element) {
                  _this.addToScene(element);
                });
                this.updateInstancesColor(proxyInstances);
                this.updateInstancesConnectionLines(proxyInstances);
                this.scene.updateMatrixWorld(true);

              case 12:
                // TODO: only update camera when cameraOptions changes
                this.cameraManager.update(cameraOptions);
                this.updateEnded();

              case 14:
              case "end":
                return _context3.stop();
            }
          }
        }, _callee3, this);
      }));

      function update(_x3, _x4, _x5, _x6, _x7) {
        return _update.apply(this, arguments);
      }

      return update;
    }()
  }, {
    key: "addToScene",
    value: function addToScene(instance) {
      var found = false;

      var _iterator8 = _createForOfIteratorHelper(this.scene.children),
          _step8;

      try {
        for (_iterator8.s(); !(_step8 = _iterator8.n()).done;) {
          var child = _step8.value;

          if (instance.instancePath && instance.instancePath === child.instancePath || child.uuid === instance.uuid) {
            found = true;
            break;
          }
        }
      } catch (err) {
        _iterator8.e(err);
      } finally {
        _iterator8.f();
      }

      if (!found) {
        this.scene.add(instance);
      }
    }
  }, {
    key: "resize",
    value: function resize() {
      if (this.width !== this.containerRef.clientWidth || this.height !== this.containerRef.clientHeight) {
        this.width = this.containerRef.clientWidth;
        this.height = this.containerRef.clientHeight;
        this.cameraManager.camera.aspect = this.width / this.height;
        this.cameraManager.camera.updateProjectionMatrix();
        this.renderer.setSize(this.width, this.height);
        this.composer.setSize(this.width, this.height); // TOFIX: this above is just an hack to trigger the ratio to be recalculated, without the line below
        // the resizing works but the image gets stretched.

        this.cameraManager.engine.controls.updateOnResize();
      }
    }
  }, {
    key: "start",
    value: function start(proxyInstances, cameraOptions, toTraverse) {
      this.resize();
      this.update(proxyInstances, cameraOptions, [], toTraverse);

      if (!this.frameId) {
        this.frameId = window.requestAnimationFrame(this.animate);
      }
    }
  }, {
    key: "requestFrame",
    value: function requestFrame() {
      var timeDif = this.lastRenderTimer.getTime() - new Date().getTime();

      if (Math.abs(timeDif) > 10) {
        this.lastRenderTimer = new Date();
        this.frameId = window.requestAnimationFrame(this.animate);
      }
    }
  }, {
    key: "animate",
    value: function animate() {
      this.controls.update();
      this.renderScene();
    }
  }, {
    key: "updateControls",
    value: function updateControls() {
      this.controls.update();
    }
  }, {
    key: "renderScene",
    value: function renderScene() {
      this.renderer.render(this.scene, this.cameraManager.getCamera());
    }
  }, {
    key: "stop",
    value: function stop() {
      cancelAnimationFrame(this.frameId);
    }
    /**
     * Returns the scene renderer
     * @returns renderer
     */

  }, {
    key: "getRenderer",
    value: function getRenderer() {
      return this.renderer;
    }
    /**
     * Returns the scene
     * @returns scene
     */

  }, {
    key: "getScene",
    value: function getScene() {
      return this.scene;
    }
    /**
     * Returns the wireframe flag
     * @returns wireframe
     */

  }, {
    key: "getWireframe",
    value: function getWireframe() {
      return this.wireframe;
    }
  }]);

  return ThreeDEngine;
}();

exports["default"] = ThreeDEngine;