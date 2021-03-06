package astar

import (
	"math"
	"testing"

	"github.com/downflux/game/map/utils"
	"github.com/downflux/game/pathing/hpf/graph"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gdpb "github.com/downflux/game/api/data_go_proto"
	mcpb "github.com/downflux/game/map/api/constants_go_proto"
	mdpb "github.com/downflux/game/map/api/data_go_proto"
	tile "github.com/downflux/game/map/map"
)

func buildTileMap(d utils.MapCoordinate, walls []utils.MapCoordinate) (*tile.Map, error) {
	wallHash := map[utils.MapCoordinate]bool{}
	for _, w := range walls {
		wallHash[w] = true
	}

	var tiles []*mdpb.Tile
	for x := int32(0); x < d.X; x++ {
		for y := int32(0); y < d.Y; y++ {
			c := utils.MC(&gdpb.Coordinate{X: x, Y: y})
			var t mcpb.TerrainType
			if _, found := wallHash[c]; found {
				t = mcpb.TerrainType_TERRAIN_TYPE_BLOCKED
			} else {
				t = mcpb.TerrainType_TERRAIN_TYPE_PLAINS
			}
			tiles = append(tiles, &mdpb.Tile{
				Coordinate:  utils.PB(c),
				TerrainType: t,
			})
		}
	}

	return tile.ImportMap(&mdpb.TileMap{
		Dimension: utils.PB(d),
		Tiles:     tiles,
		TerrainCosts: []*mdpb.TerrainCost{
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS, Cost: 1},
			{TerrainType: mcpb.TerrainType_TERRAIN_TYPE_BLOCKED, Cost: math.Inf(0)},
		},
	})
}

type aStarResult struct {
	path []*tile.Tile
	cost float64
}

func TestPath(t *testing.T) {
	singleTileMap, err := buildTileMap(utils.MC(&gdpb.Coordinate{X: 1, Y: 1}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}
	singleTileGraph, err := graph.BuildGraph(singleTileMap, &gdpb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	trivialMap, err := buildTileMap(utils.MC(&gdpb.Coordinate{X: 2, Y: 1}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}
	sourceDestinationGraph, err := graph.BuildGraph(trivialMap, &gdpb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	trivialInterClusterGraph, err := graph.BuildGraph(trivialMap, &gdpb.Coordinate{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}

	intraClusterMap, err := buildTileMap(utils.MC(&gdpb.Coordinate{X: 6, Y: 1}), nil)
	if err != nil {
		t.Fatalf("buildTileMap() = _, %v, want = _, nil", err)
	}
	trivialIntraClusterGraph, err := graph.BuildGraph(intraClusterMap, &gdpb.Coordinate{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("BuildGraph() = _, %v, want = _, nil", err)
	}
	emptyIntraClusterGraph, err := graph.BuildGraph(intraClusterMap, intraClusterMap.D)

	testConfigs := []struct {
		name string
		tm   *tile.Map
		g    *graph.Graph
		src  *gdpb.Coordinate
		dest *gdpb.Coordinate
		l    int
		want aStarResult
	}{
		{
			name: "SingleTileAStar",
			tm:   singleTileMap,
			g:    singleTileGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 0, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 0,
			},
		},
		{
			name: "SameSourceDestination",
			tm:   trivialMap,
			g:    sourceDestinationGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 0, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 0,
			},
		},
		{
			name: "TrivialInterClusterPath",
			tm:   trivialMap,
			g:    trivialInterClusterGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 1, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 1, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
		{
			name: "PartialPath",
			tm:   trivialMap,
			g:    trivialInterClusterGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 1, Y: 0},
			l:    1,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
		{
			name: "EmptyIntraClusterGraph",
			tm:   intraClusterMap,
			g:    emptyIntraClusterGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 1, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 1, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
		{
			name: "TrivialIntraClusterPath",
			tm:   intraClusterMap,
			g:    trivialIntraClusterGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 1, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 1, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 1,
			},
		},
		{
			name: "CrossIntraClusterPath",
			tm:   intraClusterMap,
			g:    trivialIntraClusterGraph,
			src:  &gdpb.Coordinate{X: 0, Y: 0},
			dest: &gdpb.Coordinate{X: 5, Y: 0},
			l:    10,
			want: aStarResult{
				path: []*tile.Tile{
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 0, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 1, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 2, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 3, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 4, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
					{
						Val: &mdpb.Tile{
							Coordinate:  &gdpb.Coordinate{X: 5, Y: 0},
							TerrainType: mcpb.TerrainType_TERRAIN_TYPE_PLAINS,
						},
					},
				},
				cost: 5,
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			path, cost, err := Path(c.tm, c.g, utils.MC(c.src), utils.MC(c.dest), c.l)
			if err != nil {
				t.Fatalf("Path() = _, _, %v, want = _, _, nil", err)
			}

			got := aStarResult{
				path: path,
				cost: cost,
			}

			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(aStarResult{}), protocmp.Transform()); diff != "" {
				t.Errorf("Path() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
