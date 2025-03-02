import { Widget } from './model';
import {BorderNode} from "./types";

/**
 * Transforms a widget configuration into a flexlayout node descriptor
 */
export function widget2Node (configuration: Widget) {
  const { id, name, component, status, panelName, enableClose = true, ...others } = configuration;
  return {
    id,
    name,
    status,
    component,
    type: "tab",
    enableRename: false,
    enableClose,
    panelName,
    ...others,
    // attr defined inside config, will also be available from within flexlayout nodes.  For example:  node.getNodeById(id).getConfig()
    config: { ...configuration } ,
  };
}

/**
 * Deep object comparison
 * @param {*} a 
 * @param {*} b 
 */
export function isEqual (a: any, b: any): boolean {
  if (a === b) {
    return true;
  }
  if (a instanceof Date && b instanceof Date) {
    return a.getTime() === b.getTime();
  }
  if (!a || !b || (typeof a !== 'object' && typeof b !== 'object')) {
    return a === b;
  }
  if (a === null || a === undefined || b === null || b === undefined) {
    return false;
  }
  if (a.prototype !== b.prototype) {
    return false;
  }
  let keys = Object.keys(a);
  if (keys.length !== Object.keys(b).length) {
    return false;
  }
  return keys.every(k => isEqual(a[k], b[k]));
}


export function findBorderToMinimizeTo(borderNodes: BorderNode[]): BorderNode | null {
  for (let borderNode of borderNodes) {
    if (borderNode.getConfig()?.isMinimizedPanel) {
      return borderNode;
    }
  }
  return null;
}

export function getWidget(store, id){
  return store.getState().widgets[id];
}