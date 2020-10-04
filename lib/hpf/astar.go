package astar

import (
	"math"

	"github.com/minkezhang/rts-pathing/lib/hpf/cluster"
	"github.com/minkezhang/rts-pathing/lib/hpf/graph"
	"github.com/minkezhang/rts-pathing/lib/hpf/graphastar"
	"github.com/minkezhang/rts-pathing/lib/hpf/tile"
	"github.com/minkezhang/rts-pathing/lib/hpf/tileastar"
	"github.com/minkezhang/rts-pathing/lib/hpf/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

func clusterBoundedTilePath(tm *tile.Map, g *graph.Graph, src, dest utils.MapCoordinate) ([]*tile.Tile, float64, error) {
	c1, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, src)
	if err != nil {
		return nil, 0, err
	}
	c2, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, dest)
	if err != nil {
		return nil, 0, err
	}

	if c1 != c2 {
		return nil, 0, status.Error(codes.FailedPrecondition, "input source and destination nodes do not exist in the same cluster")
	}

	tileBoundary, err := cluster.TileBoundary(g.NodeMap.ClusterMap, c1)
	if err != nil {
		return nil, 0, err
	}
	tileDimension, err := cluster.TileDimension(g.NodeMap.ClusterMap, c1)
	if err != nil {
		return nil, 0, err
	}

	return tileastar.Path(tm, src, dest, utils.PB(tileBoundary), utils.PB(tileDimension))
}

func Path(tm *tile.Map, g *graph.Graph, src, dest utils.MapCoordinate, l int) ([]*tile.Tile, float64, error) {
	if l < 0 {
		return nil, 0, status.Error(codes.FailedPrecondition, "cannot specify a negative path length")
	}

	p, c, _ := clusterBoundedTilePath(tm, g, src, dest)
	if p != nil {
		return p, c, nil
	}

	srcID, err := graph.InsertEphemeralNode(tm, g, src)
	if err != nil {
		return nil, 0, err
	}
	defer graph.RemoveEphemeralNode(g, src, srcID)

	destID, err := graph.InsertEphemeralNode(tm, g, dest)
	if err != nil {
		return nil, 0, err
	}
	defer graph.RemoveEphemeralNode(g, dest, destID)

	nPath, cost, err := graphastar.Path(tm, g, src, dest)
	if err != nil {
		return nil, 0, err
	}
	if math.IsInf(cost, 0) {
		return nil, cost, nil
	}

	var path []*tile.Tile
	for i, n1 := range nPath {
		t1 := utils.MC(n1.GetTileCoordinate())
		c1, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t1)
		if err != nil {
			return nil, 0, err
		}

		// The implementations of astar returns both the source and
		// destination Tile instances. If we naively append all
		// segmented paths together, we will get duplicates as the
		// destination of one segment is the source of the following
		// segment. In order to avoid these duplicates, we strip the
		// source from all path results; this code block initializes
		// the returned path by adding the global source Tile.
		if i == 0 {
			path = append(path, tm.TileFromCoordinate(utils.PB(t1)))
		}
		// Last element in an AbstractNode list do not have a
		// corresponding "target" to move to.
		if i+1 == len(nPath) {
			break
		}

		n2 := nPath[i+1]
		t2 := utils.MC(n2.GetTileCoordinate())
		c2, err := cluster.ClusterCoordinateFromTileCoordinate(g.NodeMap.ClusterMap, t2)
		if err != nil {
			return nil, 0, err
		}

		var p []*tile.Tile
		if c1 == c2 {
			p, _, err := clusterBoundedTilePath(tm, g, t1, t2)
			if err != nil {
				return nil, 0, err
			}
			_, p = p[0], p[1:]
		} else {
			// Inter-cluster nodes are always immediately adjacent
			// and unblocked.
			p = append(p, tm.TileFromCoordinate(utils.PB(t2)))
		}

		for _, n := range p {
			if l == 0 || len(p) < l {
				path = append(path, n)
			}
		}
	}

	return path, cost, err
}
