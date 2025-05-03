
// Tunable parameters
const zoomSensitivityMouse = 0.001; // Adjust for zoom speed
const zoomSensitivityTouchpad = 0.01; // Adjust for zoom speed
const dragMultiplier = 2.5; // Adjust for drag speed

const root = document.querySelector('svg#root');
const map = document.querySelector('svg#map');
const legend = document.querySelector('g#legend');
const timeline = document.querySelector('svg#timeline');

function updatePositions() {
    const {width: rootWidth, height: rootHeight} = root.getBoundingClientRect(); 
    const legendHeight = legend.getBoundingClientRect().height; 
    let mapHeight = map.getBoundingClientRect().height; 
    let timelineHeight = timeline.getBoundingClientRect().height;
    mapHeight = rootHeight - legendHeight - timelineHeight
    map.setAttribute('height', mapHeight)
    timeline.setAttribute('y', rootHeight - timelineHeight-20)
    timeline.setAttribute('width', rootWidth - 40)
};

updatePositions();

window.addEventListener("resize", updatePositions)

// map zooming/panning
const [_minX, _minY, maxX, maxY] = map.getAttribute('viewBox').split(' ').map(Number);
let startX, startY, initialDragViewBox;

function clamp(x, max) {
    if (0 > x) return 0;
    if (x > max) return max;
    return x
}

function mapDragStart(event) {
    map.style.cursor = 'grabbing';
    startX = event.clientX;
    startY = event.clientY;
    initialDragViewBox = map.getAttribute('viewBox').split(' ').map(Number);
    map.addEventListener('mousemove', mapDragMove);

};

function mapDragMove(event) {
    const deltaX = event.clientX - startX;
    const deltaY = event.clientY - startY;

    // Adjust the viewBox based on the mouse movement
    const newMinX = clamp(initialDragViewBox[0] - deltaX*dragMultiplier, maxX);
    const newMinY = clamp(initialDragViewBox[1] - deltaY*dragMultiplier, maxY);

    map.setAttribute('viewBox', `${newMinX} ${newMinY} ${initialDragViewBox[2]} ${initialDragViewBox[3]}`);
};

function mapDragStop() {
    map.removeEventListener('mousemove', mapDragMove);
    map.style.cursor = 'auto';
};

function zoomSensitivity(event) {
    // Check if deltaX or deltaY are very small and non-zero, which is typical for touchpads
    if (1 > Math.abs(event.deltaY) || 1 > Math.abs(event.deltaX)) {
        return zoomSensitivityTouchpad;
    } else {
        return zoomSensitivityMouse;
    }
}

function mapWheelZoom(event) {
    event.preventDefault(); // Prevent default scrolling

    const delta = -event.deltaY * zoomSensitivity(event); // Normalize scroll direction and sensitivity

    let [minX, minY, width, height] = map.getAttribute('viewBox').split(' ').map(Number);

    // Calculate the point to zoom around (mouse position relative to SVG)
    const rect = map.getBoundingClientRect();
    const mouseX = event.clientX - rect.left;
    const mouseY = event.clientY - rect.top;

    // Calculate the focal point in the SVG's user coordinates
    const focusX = minX + width * (mouseX / rect.width);
    const focusY = minY + height * (mouseY / rect.height);

    // Adjust width and height based on zoom delta
    const newWidth = clamp(width * (1 + delta), maxX);
    const newHeight = clamp(height * (1 + delta), maxY);

    // Ensure width and height don't become zero or negative
    if (0 == newWidth || 0 == newHeight) return;

    // Adjust minX and minY to zoom around the focal point
    const newMinX = clamp(focusX - (mouseX / rect.width) * newWidth, maxX);
    const newMinY = clamp(focusY - (mouseY / rect.height) * newHeight, maxY);

    map.setAttribute('viewBox', `${newMinX} ${newMinY} ${newWidth} ${newHeight}`);
};

map.addEventListener('mousedown', mapDragStart);
map.addEventListener('mouseup', mapDragStop);
map.addEventListener('mouseleave', mapDragStop);
map.addEventListener('wheel', mapWheelZoom , { passive: false }); // passive: false is needed to preventDefault

// map hover handling
// * highlight the corresponding timeline segment

function getMapSegmentId(event) {
    const target = event.target
    if(!target) return undefined
    if(target.tagName != 'line') return undefined
    const id = event.target.parentNode.getAttribute("id")
    return id
}

function mapSegmentHoverStart(event) {
    const id = getMapSegmentId(event)
    if (!id) return
    const selector = `rect#${id}`
    const timelineSegment = timeline.querySelector(selector)
    if (!timelineSegment) return
    timelineSegment.classList.add("timeline-segment-rect-hovered")
}

function mapSegmentHoverStop(event) {
    const id = getMapSegmentId(event)
    if (!id) return
    const selector = `rect#${id}`
    const timelineSegment = timeline.querySelector(selector)
    if (!timelineSegment) return
    timelineSegment.classList.remove("timeline-segment-rect-hovered")
}

map.addEventListener('mouseover', mapSegmentHoverStart)
map.addEventListener('mouseout', mapSegmentHoverStop)

// timeline hover handling
// * highlight corresponding map segment
function getTimelineSegmentId(event) {
    const target = event.target
    if(!target) return undefined
    if(target.tagName != 'rect' && target.tagName != 'polygon') return undefined
    const id = event.target.getAttribute("id")
    return id
}

function timelineSegmentHoverStart(event) {
    const id = getTimelineSegmentId(event)
    if (!id) return
    const selector = `g#${id}.segment`
    const mapSegment = map.querySelector(selector)
    if (!mapSegment) return
    mapSegment.classList.add("segment-hovered")
}

function timelineSegmentHoverStop(event) {
    const id = getTimelineSegmentId(event)
    if (!id) return
    const mapSegment = map.querySelector(`g#${id}.segment`)
    if (!mapSegment) return
    mapSegment.classList.remove("segment-hovered")
}

timeline.addEventListener('mouseover', timelineSegmentHoverStart)
timeline.addEventListener('mouseout', timelineSegmentHoverStop)

// timeline period selection
// * show box around selected segments
// * hide everything but selected segments in the map
// * narrow timeline viewbox to selected segments

// Range represents a segment ID range as a continuous interval
// (as opposed to a discrete set of IDs)
class Range {
    constructor(min, max) {
        this.min = Range.toNum(min)
        this.max = Range.toNum(max) ?? Range.toNum(min)
    }

    static toNum(id) {
        if (!id) return undefined;
        return parseInt(id.slice(1))
    }

    add(id) {
        const val = Range.toNum(id)
        if (!this.min || val < this.min) this.min = val
        if (!this.max || val > this.max) this.max = val
    }

    has(id) {
        const val = Range.toNum(id)
        return this.min && this.min <= val && this.max && val <= this.max
    }
}

let selectionBox = null; // tracks the visual timeline selection
let selectedSegments = null; // tracks the selected segment range

function timelineSelectStart(event) {
    const id = getTimelineSegmentId(event)
    if (!id) return
    // initialize the selecteSegments range and the selection box
    selectedSegments = new Range(id);
    const segmentBox = timeline.querySelector(`rect#${id}`)
    selectionBox = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
    selectionBox.classList.add('timeline-selection-box')
    selectionBox.setAttribute('id', 'selection-box');
    selectionBox.setAttribute('x', segmentBox.getAttribute("x"));
    selectionBox.setAttribute('y', 0);
    selectionBox.setAttribute('width', segmentBox.getAttribute("width"));
    selectionBox.setAttribute('height', "100%");
    timeline.appendChild(selectionBox);
    // start tracking the selection
    timeline.addEventListener("mousemove", timelineSelectMove)
}

function timelineSelectMove(event) {
    // add the target segment under cursor to the selectedSegments range.
    const id = getTimelineSegmentId(event)
    if (!id || id == selectionBox.getAttribute('id') || selectedSegments.has(id)) return
    selectedSegments.add(id);
    // expand the selectionBox to include the box of the added segment
    const segmentBox = timeline.querySelector(`rect#${id}`)
    const x = parseInt(selectionBox.getAttribute('x'))
    let width = parseInt(selectionBox.getAttribute('width'))
    width += parseInt(segmentBox.getAttribute('width')) 
    selectionBox.setAttribute('width', width);
}

function timelineSelectStop(event) {
    if (!selectionBox) return
    // stop tracking selection
    timeline.removeEventListener("mousemove", timelineSelectMove)
    // zoom timeline viewport to show only the range of the selection box
    let viewBox = timeline.getAttribute('viewBox')
    const [_x, y, _width, height] = viewBox.split(' ').map(Number);
    viewBox = `${selectionBox.getAttribute('x')} ${y} ${selectionBox.getAttribute('width')} ${height}`;
    timeline.setAttribute('viewBox', viewBox)
    // drop the selection box
    selectionBox.remove()
    selectionBox = null
    // go over all map segments and hide the ones that aren't in selectedSegments range
    for (const elem of map.children) {
        if (elem.tagName != 'g' || selectedSegments.has(elem.getAttribute('id'))) continue;
        elem.style.visibility = 'hidden';
    }
    selectedSegments = null
}

timeline.addEventListener("mousedown", timelineSelectStart)
timeline.addEventListener("mouseup", timelineSelectStop)