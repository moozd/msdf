# MSDF Approach 2: Direct Multi-Channel Distance Field Construction

## Overview

Approach 2 from Viktor Chlumský's thesis "Shape Decomposition for Multi-channel Distance Fields" describes the **direct method** for constructing multi-channel signed distance fields directly from vector representation without intermediate raster decomposition. This approach is superior in both performance and quality compared to approach 1.

## Key Principles

### 1. Multi-Channel Distance Fields
- Uses **3 channels (RGB)** to store different distance field information
- Each channel represents distance to different subset of shape edges
- Final rendering uses **median-of-three** reconstruction: `median(R, G, B)`
- This preserves sharp corners that single-channel SDF would blur

### 2. Edge Coloring Strategy
The core innovation is **intelligent edge coloring**:
- Each edge gets assigned one or more colors from {RED, GREEN, BLUE}
- Adjacent edges must share at least one color
- Goal: ensure every corner has edges with different color combinations
- This prevents distance field conflicts at sharp corners

## Step-by-Step Implementation

### Step 1: Vector Shape Analysis
1. **Parse font glyph** into vector contours
2. **Identify edges**: lines, quadratic Bézier, cubic Bézier curves
3. **Determine winding direction** (clockwise/counter-clockwise) for each contour
4. **Build edge connectivity graph** to understand corner relationships

### Step 2: Edge Coloring Algorithm
```
For each contour:
  edges = get_edges(contour)
  
  if edges.length == 1:
    edge.color = WHITE  // All three channels
  else:
    // Start with two colors
    current_color = RED | BLUE
    
    for each edge in edges:
      edge.color = current_color
      
      // Rotate colors to ensure adjacent edges share colors
      // but corners have different color combinations
      if current_color has (RED AND GREEN):
        current_color = GREEN | BLUE
      else if current_color has (GREEN AND BLUE):
        current_color = RED | GREEN  
      else:
        current_color = RED | GREEN
```

### Step 3: Distance Field Computation
For each pixel (x, y) and each channel (R, G, B):

1. **Find closest edge** with that channel color:
   ```
   min_distance = infinity
   closest_point = null
   
   for each edge with channel_color:
     distance, point = compute_distance_to_edge(x, y, edge)
     if distance < min_distance:
       min_distance = distance
       closest_point = point
   ```

2. **Compute accurate distance** based on curve type:
   - **Lines**: Point-to-line segment distance
   - **Quadratic Bézier**: Solve quadratic equation for closest point
   - **Cubic Bézier**: Numerical method or polynomial root finding

3. **Apply sign based on winding**:
   ```
   // Determine if point is inside or outside shape
   edge_direction = edge.tangent_at_closest_point
   point_to_closest = vector(closest_point, pixel)
   
   side = sign(cross_product(edge_direction, point_to_closest))
   signed_distance = side * winding_direction * distance
   ```

### Step 4: Distance Field Normalization
```
distance_range = 0.5 * pixel_size  // Usually 0.5-2 pixels
normalized = (signed_distance / distance_range) + 0.5
clamped = clamp(normalized, 0.0, 1.0)
channel_value = clamped * 255
```

### Step 5: Multi-Channel Reconstruction
The final rendering uses **median filtering**:
```glsl
// In fragment shader
vec3 msd = texture(msdf_texture, uv).rgb;
float sd = median(msd.r, msd.g, msd.b);
float alpha = clamp(sd + 0.5, 0.0, 1.0);
```

## Mathematical Foundation

### Bézier Distance Computation

**Quadratic Bézier**: P(t) = (1-t)²P₀ + 2t(1-t)P₁ + t²P₂
- Find t where dP/dt ⊥ (P(t) - query_point)
- Solve: dot(P'(t), P(t) - Q) = 0

**Cubic Bézier**: P(t) = (1-t)³P₀ + 3t(1-t)²P₁ + 3t²(1-t)P₂ + t³P₃
- More complex: solve cubic equation or use numerical methods
- Newton-Raphson iteration for root finding

### Winding Number Computation
```
winding = 0
for each edge in contour:
  if edge crosses horizontal ray from point:
    if edge goes upward: winding += 1
    else: winding -= 1

inside = (winding != 0)
```

## Advantages of Approach 2

1. **No intermediate rasterization** - works directly with vectors
2. **Better quality** - preserves mathematical precision
3. **Faster** - eliminates raster decomposition step
4. **Scalable** - works at any resolution
5. **Memory efficient** - no temporary raster buffers

## Implementation Checklist

- [ ] Parse font glyphs into vector contours
- [ ] Implement proper edge coloring algorithm
- [ ] Add accurate Bézier curve distance computation
- [ ] Handle winding direction correctly
- [ ] Implement signed distance calculation
- [ ] Add distance field normalization
- [ ] Create proper multi-channel output
- [ ] Test with median-of-three reconstruction shader

## Common Pitfalls

1. **Incorrect edge coloring** - leads to color bleeding at corners
2. **Incomplete curve distance** - causes artifacts in curved sections
3. **Wrong winding calculation** - inverts inside/outside regions
4. **Poor normalization** - results in aliasing or blurring
5. **Missing median reconstruction** - defeats purpose of multi-channel approach

## References

- Chlumský, Viktor. "Shape Decomposition for Multi-channel Distance Fields" (Master's Thesis)
- Chapter 4: "Direct multi-channel distance field construction"
- Section 4.2: "Edge coloring"
- Section 4.3: "Distance field computation"