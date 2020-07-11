package entrance

import (
	"testing"

	rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/cluster"
	// "github.com/cripplet/rts-pathing/lib/hpf/tile"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestbuildClusterEdgeCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		c    *cluster.Cluster
		d    rtscpb.Direction
	}{
		{name: "NullClusterTest", c: cluster.ImportCluster(&rtsspb.Cluster{}), d: rtscpb.Direction_DIRECTION_NORTH},
		{name: "NullXDimensionClusterTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 0, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "NullYDimensionClusterTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 0},
			}), d: rtscpb.Direction_DIRECTION_NORTH,
		},
		{name: "InvalidDirectionTest", c: cluster.ImportCluster(
			&rtsspb.Cluster{
				TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
				TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
			}), d: rtscpb.Direction_DIRECTION_UNKNOWN,
		},
	}
	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildClusterEdgeCoordinateSlice(c.c, c.d); err == nil {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want a non-nil error", got, err)
			}
		})
	}
}

func TestBuildClusterEdgeCoordinateSlice(t *testing.T) {
	trivialCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 1},
	})
	smallCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 0, Y: 0},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	embeddedCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 2, Y: 2},
	})
	rectangularCluster := cluster.ImportCluster(&rtsspb.Cluster{
		TileBoundary:  &rtsspb.Coordinate{X: 1, Y: 1},
		TileDimension: &rtsspb.Coordinate{X: 1, Y: 2},
	})
	testConfigs := []struct {
		name string
		c    *cluster.Cluster
		d    rtscpb.Direction
		want *rtsspb.CoordinateSlice
	}{
		{name: "TrivialClusterNorthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterSouthTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterEastTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "TrivialClusterWestTest", c: trivialCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}},
		{name: "SmallClusterNorthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2}},
		{name: "SmallClusterSouthTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "SmallClusterEastTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2}},
		{name: "SmallClusterWestTest", c: smallCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2}},
		{name: "EmbeddedClusterNorthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 2}},
		{name: "EmbeddedClusterSouthTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "EmbeddedClusterEastTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 1}, Length: 2}},
		{name: "EmbeddedClusterWestTest", c: embeddedCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "RectangularClusterNorthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_NORTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 1}},
		{name: "RectangularClusterSouthTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_SOUTH, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 1}},
		{name: "RectangularClusterEastTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_EAST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
		{name: "RectangularClusterWestTest", c: rectangularCluster, d: rtscpb.Direction_DIRECTION_WEST, want: &rtsspb.CoordinateSlice{
			Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2}},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildClusterEdgeCoordinateSlice(c.c, c.d); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildClusterEdgeCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *rtsspb.CoordinateSlice
		o    int32
		want *rtsspb.Coordinate
	}{
		{
			name: "SingleTileSliceHorizontal",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			o:    0,
			want: &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name: "SingleTileSliceVertical",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			o:    0,
			want: &rtsspb.Coordinate{X: 0, Y: 0},
		},
		{
			name: "MultiTileTileSliceHorizontal",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			o:    1,
			want: &rtsspb.Coordinate{X: 2, Y: 1},
		},
		{
			name: "MultiTileTileSliceVertical",
			s:    &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
			o:    1,
			want: &rtsspb.Coordinate{X: 1, Y: 2},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildCoordinateWithCoordinateSlice(c.s, c.o); err != nil || !proto.Equal(got, c.want) {
				t.Errorf("buildCoordinateWithCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildCoordinateWithCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name string
		s    *rtsspb.CoordinateSlice
		o    int32
	}{
		{name: "NullTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 0}, o: 0},
		{name: "OutOfBoundsTileSliceBefore", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: -1},
		{name: "OutOfBoundsTileSliceAfter", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: 2},
		{name: "InvalidOrientationTileSlice", s: &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_UNKNOWN, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1}, o: 0},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if _, err := buildCoordinateWithCoordinateSlice(c.s, c.o); err == nil {
				t.Error("buildCoordinateWithCoordinateSlice() = nil, want a non-nil error")
			}
		})
	}
}

func TestBuildTransitionsFromOpenCoordinateSlice(t *testing.T) {
	testConfigs := []struct {
		name   string
		s1, s2 *rtsspb.CoordinateSlice
		want   []*rtsspb.Transition
	}{
		{
			name: "SingleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 1},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 0, Y: 1}},
			},
		},
		{
			name: "SingleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 1},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 0}},
			},
		},
		{
			name: "DoubleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "DoubleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 2},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "TripleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 3},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "TripleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 3},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 3},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 1}},
			},
		},
		{
			name: "QuadrupleWidthEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 0, Y: 1}},
				{C1: &rtsspb.Coordinate{X: 3, Y: 0}, C2: &rtsspb.Coordinate{X: 3, Y: 1}},
			},
		},
		{
			name: "QuadrupleWidthEntranceVertical",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 0}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 0, Y: 0}, C2: &rtsspb.Coordinate{X: 1, Y: 0}},
				{C1: &rtsspb.Coordinate{X: 0, Y: 3}, C2: &rtsspb.Coordinate{X: 1, Y: 3}},
			},
		},
		{
			name: "QuadrupleWidthEmbeddedEntranceHorizontal",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 4},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 2}, Length: 4},
			want: []*rtsspb.Transition{
				{C1: &rtsspb.Coordinate{X: 1, Y: 1}, C2: &rtsspb.Coordinate{X: 1, Y: 2}},
				{C1: &rtsspb.Coordinate{X: 4, Y: 1}, C2: &rtsspb.Coordinate{X: 4, Y: 2}},
			},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if got, err := buildTransitionsFromOpenCoordinateSlice(c.s1, c.s2); err != nil || !cmp.Equal(got, c.want, protocmp.Transform()) {
				t.Errorf("buildTransitionsFromOpenCoordinateSlice() = %v, %v, want = %v, nil", got, err, c.want)
			}
		})
	}
}

func TestBuildTransitionsFromOpenCoordinateSliceError(t *testing.T) {
	testConfigs := []struct {
		name   string
		s1, s2 *rtsspb.CoordinateSlice
	}{
		{
			name: "MismatchedLengths",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 1}, Length: 2},
		},
		{
			name: "MismatchedOrientations",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
		},
		{
			name: "NonAdjacentHorizontalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 2}, Length: 1},
		},
		{
			name: "NonAdjacentVerticalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 1},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 2, Y: 0}, Length: 1},
		},
		{
			name: "NonAlignedHorizontalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_HORIZONTAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
		},
		{
			name: "NonAlignedVerticalSlice",
			s1:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 0, Y: 0}, Length: 2},
			s2:   &rtsspb.CoordinateSlice{Orientation: rtscpb.Orientation_ORIENTATION_VERTICAL, Start: &rtsspb.Coordinate{X: 1, Y: 1}, Length: 2},
		},
	}

	for _, c := range testConfigs {
		t.Run(c.name, func(t *testing.T) {
			if _, err := buildTransitionsFromOpenCoordinateSlice(c.s1, c.s2); err == nil {
				t.Error("buildTransitionsFromOpenCoordinateSlice() = nil, want a non-nil error")
			}
		})
	}
}
