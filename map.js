
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
    console.log(rootHeight, legendHeight, mapHeight, timelineHeight)
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

function dragStart(event) {
    map.style.cursor = 'grabbing';
    startX = event.clientX;
    startY = event.clientY;
    initialDragViewBox = map.getAttribute('viewBox').split(' ').map(Number);
    map.addEventListener('mousemove', dragMove);

};

function dragMove(event) {
    const deltaX = event.clientX - startX;
    const deltaY = event.clientY - startY;

    // Adjust the viewBox based on the mouse movement
    const newMinX = clamp(initialDragViewBox[0] - deltaX*dragMultiplier, maxX);
    const newMinY = clamp(initialDragViewBox[1] - deltaY*dragMultiplier, maxY);

    map.setAttribute('viewBox', `${newMinX} ${newMinY} ${initialDragViewBox[2]} ${initialDragViewBox[3]}`);
};

function dragStop() {
    map.removeEventListener('mousemove', dragMove);
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

function wheelZoom(event) {
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

map.addEventListener('mousedown', dragStart);
map.addEventListener('mouseup', dragStop);
map.addEventListener('mouseleave', dragStop);
map.addEventListener('wheel', wheelZoom , { passive: false }); // passive: false is needed to preventDefault

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